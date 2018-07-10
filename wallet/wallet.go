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
const addressChecksumLen = 4

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates and returns a Wallet
func NewWallet() (*Wallet, error) {
	private, public, err := newKeyPair()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	wallet := Wallet{private, public}

	return &wallet, nil
}

func (w Wallet) PubKeyHash() ([]byte, error) {
	return HashPubKey(w.PublicKey)
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() ([]byte, error) {
	pubKeyHash, err := HashPubKey(w.PublicKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Println(hex.EncodeToString(w.PublicKey))
	log.Println(pubKeyHash)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address, err
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) ([]byte, error) {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160, nil
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

func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Println(err.Error())
		return ecdsa.PrivateKey{}, nil, err
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey, nil
}

func (w Wallet) CreateTx(chain *core.Blockchain, to []byte, amount uint32, utxo *core.UTXOSet) (*types.Transaction, error) {
	var inputs []*types.TxIn
	var outputs []*types.TxOut

	pubKeyHash, err := HashPubKey(w.PublicKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	acc, validOutputs := utxo.FindTxOutsOfAmount(pubKeyHash, amount)

	if acc < amount {
		err := errors.New("ERROR: Not enough funds")
		log.Println(err.Error())
		return nil, err
	}

	log.Println(len(validOutputs), validOutputs)
	// Build a list of inputs
	for txid, outs := range validOutputs {
		for _, out := range outs {
			input := &types.TxIn{common.ToHash32([]byte(txid)), int64(out), nil, pubKeyHash}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, &types.TxOut{amount, to})
	if acc > amount {
		outputs = append(outputs, &types.TxOut{acc - amount, pubKeyHash})
	}

	log.Println(len(inputs), *inputs[0])
	log.Println(len(outputs), *outputs[0], *outputs[1])
	tx := &types.Transaction{LockTime: uint32(time.Now().Unix()), TxIn: inputs, TxOut: outputs}
	tx.TxHash, err = tx.Hash()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	chain.SignTransaction(tx, w.PrivateKey)

	return tx, nil
}
