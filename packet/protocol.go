package packet

import (
	"bytes"
	"encoding/gob"
)

type Body map[string]interface{}

// implement error interface
func (b Body) Error() string {
	return b["msg"].(string)
}

func (b Body) Decode(bts []byte) error {
	var buf bytes.Buffer
	if _, err := buf.Read(bts); err != nil {
		return err
	}
	gob.Register(Body{})
	decoder := gob.NewDecoder(&buf)

	if err := decoder.Decode(b); err != nil {
		return err
	}
	return nil
}
