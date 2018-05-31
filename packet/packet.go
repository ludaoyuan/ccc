package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io/ioutil"
)

// const (
// 	RP = [8]byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55}
// 	WP = [8]byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff}
// )

type Body map[string]interface{}

func (b Body) Error() string {
	return b["msg"].(string)
}

type Packet struct {
	Header	[512]byte
	Body	[]byte
}

func (p Packet) Validation(s [512]byte, n int, body *[]byte) (bool, Body) {
	var buf bytes.Buffer

	gob.Register(Body {})

	bm := make(Body)
	encoder := gob.NewEncoder(&buf)

	bm["code"] = uint16(0x0000)
	bm["msg"] = "success"
	ok := true

	if n != 512 || !bytes.Equal(p.Header[:8], s[:8]) || !bytes.Equal(p.Header[504:], s[504:]) {
		bm["code"] = uint16(0x0001)
		bm["msg"] = "not a inc packet"
		ok = false
	}

	err := encoder.Encode(bm)
	if err != nil {
		bm["code"] = uint(0x1000) // 0x1000 原生错误
		bm["msg"] = err.Error()
		ok = false
		return ok, bm
	}

	*body, err = ioutil.ReadAll(&buf)
	if err != nil {
		bm["code"] = uint(0x1000) // 0x1000 原生错误
		bm["msg"] = err.Error()
		ok = false
		return ok, bm
	}

	return ok, nil
}

func PacketsCreation() (*Packet, *Packet) {
	rPacket := &Packet{[512]byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55}, []byte{}}
	copy(rPacket.Header[504:], []byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55})
	wPacket := &Packet{[512]byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff}, []byte{}}
	copy(wPacket.Header[504:], []byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff})

	return rPacket, wPacket
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

func (p *Packet) GetLength() int {
	return int(binary.BigEndian.Uint64(p.Header[8:16]))
}

func (p *Packet) SetLength() {
	binary.BigEndian.PutUint64(p.Header[8:16], uint64(len(p.Body)))
}

/*
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
*/
