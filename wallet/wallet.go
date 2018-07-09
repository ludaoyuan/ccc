package wallet

import (
	"bytes"
	"common"
	"core"
	"core/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "./data/wallet.dat"
const addressChecksumLen = 4

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// 创建交易
func (w Wallet) CreateTx(chain *core.Blockchain, to []byte, amount uint32, utxos *core.UTXOSet) (*types.Transaction, error) {
	var ins []*types.TxIn
	var outs []*types.TxOut

	pubKeyHash := HashPubKey(w.PublicKey)
	acc, validOuts := utxos.FindTxOutsOfAmount(pubKeyHash, amount)

	if acc < amount {
		err := errors.New("ERROR: Not enough funds")
		log.Println(err.Error())
		return nil, err
	}

	for txid, outs := range validOuts {
		idbytes, err := hex.DecodeString(txid)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		for _, out := range outs {
			in := &types.TxIn{common.ToHash(idbytes), int64(out), [32]byte{}, common.ToHash(w.PublicKey)}
			ins = append(ins, in)
		}
	}

	from := w.GetAddress()
	outs = append(outs, &types.TxOut{amount, common.ToHash(to)})
	if acc > amount {
		outs = append(outs, &types.TxOut{acc - amount, common.ToHash(from)}) // a change
	}

	tx := types.Transaction{[32]byte{}, uint32(time.Now().Unix()), ins, outs}
	var err error
	tx.TxHash, err = tx.Hash()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	chain.SignTransaction(&tx, w.PrivateKey)

	return &tx, nil
}
