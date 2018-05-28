package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
)

// const (
// 	RP = [8]byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55}
// 	WP = [8]byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff}
// )

type Packet struct {
	Header [512]byte
	Body   []byte
}

func (p Packet) Validation(s [512]byte, n int, body *[]byte) bool {
	log.Println(s[504:])
	log.Println(p.Header[504:])
	if n != 512 || !bytes.Equal(p.Header[:8], s[:8]) || !bytes.Equal(p.Header[504:], s[504:]) {
		*body = []byte(`"msg":"not a nnc packet"`)
		return false
	}

	return true
}

func RPacket() *Packet {
	rPacket := &Packet{[512]byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55}, []byte{}}
	copy(rPacket.Header[504:], []byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55})

	rPacket.Body = make([]byte, 0)

	return rPacket
}

func WPacket() *Packet {
	wPacket := &Packet{[512]byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff}, []byte{}}
	copy(wPacket.Header[504:], []byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff})

	wPacket.Body = make([]byte, 0)

	return wPacket
}

func (p *Packet) ResetRPacket() {
	p.Body = p.Body[:0]
	copy(p.Header[:], []byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55}[:])
	copy(p.Header[504:], []byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55})
}

func (p *Packet) ResetWPacket() {
	p.Body = p.Body[:0]
	copy(p.Header[:], []byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff}[:])
	copy(p.Header[504:], []byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff})
}

func (p Packet) Length() int {
	return int(binary.BigEndian.Uint64(p.Header[8:16]))
}

func (p *Packet) SetLength(length uint64) {
	binary.BigEndian.PutUint64(p.Header[8:16], length)
}

func (p *Packet) Decode() (*Msg, error) {
	var buf bytes.Buffer
	decoder := gob.NewDecoder(&buf)
	msg := Msg{}
	if err := decoder.Decode(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (p *Packet) Err() error {
	msg, err := p.Decode()
	if err != nil {
		return err
	}
	log.Println(msg)
	return msg.ErrMsg
}
