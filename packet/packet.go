package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io/ioutil"
)

type Packet struct {
	Header [512]byte
	Body   []byte
}

func (p Packet) Validation(s [512]byte, n int, body *[]byte) Body {
	var buf bytes.Buffer

	gob.Register(Body{})

	bm := make(Body)
	encoder := gob.NewEncoder(&buf)

	if n != 512 || !bytes.Equal(p.Header[:8], s[504:]) || !bytes.Equal(p.Header[504:], s[:8]) {
		bm["code"] = uint16(0x0001)
		bm["msg"] = "not a inc packet"
		return bm
	}

	err := encoder.Encode(bm)
	if err != nil {
		bm["code"] = uint(0x1000) // 0x1000 原生错误
		bm["msg"] = err.Error()
		return bm
	}

	*body, err = ioutil.ReadAll(&buf)
	if err != nil {
		bm["code"] = uint(0x1000) // 0x1000 原生错误
		bm["msg"] = err.Error()
		return bm
	}

	return nil
}

func PacketsCreation() (*Packet, *Packet) {
	rPacket := &Packet{[512]byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55}, []byte{}}
	copy(rPacket.Header[504:], []byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff})
	wPacket := &Packet{[512]byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff}, []byte{}}
	copy(wPacket.Header[504:], []byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55})

	return rPacket, wPacket
}

func (p *Packet) ResetRPacket() {
	p.Body = p.Body[:0]
	copy(p.Header[:], []byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55}[:])
	copy(p.Header[504:], []byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff})
}

func (p *Packet) ResetWPacket() {
	p.Body = p.Body[:0]
	copy(p.Header[:], []byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff}[:])
	copy(p.Header[504:], []byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55})
}

func (p *Packet) GetLength() int {
	return int(binary.BigEndian.Uint64(p.Header[8:16]))
}

func (p *Packet) SetLength() {
	binary.BigEndian.PutUint64(p.Header[8:16], uint64(len(p.Body)))
}
