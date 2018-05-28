package packet

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// error codes used by this  protocol scheme
const (
	ErrMsgTooLong = iota
	ErrDecode
	ErrWrite
	ErrInvalidMsgCode
	ErrInvalidMsgType
	ErrHandshake
	ErrNoHandler
	ErrHandler
)

// error description strings associated with the codes
var errorToString = map[int]string{
	ErrMsgTooLong:     "Message too long",
	ErrDecode:         "Invalid message (RLP error)",
	ErrWrite:          "Error sending message",
	ErrInvalidMsgCode: "Invalid message code",
	ErrInvalidMsgType: "Invalid message type",
	ErrHandshake:      "Handshake error",
	ErrNoHandler:      "No handler registered error",
	ErrHandler:        "Message handler error",
}

type ErrMsg struct {
	ErrCode uint64
	ErrInfo []byte
}

// Error implement of Error interface
func (err ErrMsg) Error() string {
	return fmt.Sprintf("Code: %d, Info: %s", err.ErrCode, err.ErrInfo)
}

func NewErrMsg(errInfo []byte) *ErrMsg {
	return &ErrMsg{
		ErrInfo: errInfo,
	}
}

type Msg struct {
	ErrMsg  *ErrMsg
	Message interface{}
}

func (m *Msg) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Msg) Err() error {
	return m.ErrMsg
}
