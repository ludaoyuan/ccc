package packet

import (
	"bytes"
	"log"
)

type Packet struct {
	Header	[512]byte
	Body	[]byte
}

func (p Packet)Validation(s [512]byte, n int, body *[]byte) bool {
	log.Println(s[504:])
	log.Println(p.Header[504:])
	if n != 512 || !bytes.Equal(p.Header[:8], s[:8]) || !bytes.Equal(p.Header[504:], s[504:]) {
		*body = []byte(`"msg":"not a nnc packet"`)
		return false
	}

	return true
}
