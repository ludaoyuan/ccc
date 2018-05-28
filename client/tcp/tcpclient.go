package client

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"log"
	"net"
	"packet"
)

// Client is the client for communicating with server
type Client struct {
	conn   net.Conn
	Buffer *packet.Packet
}

// TODO: 日志模块，选择一个等级日志模块
/*
Run
ial
Read
Write
Buffer
*/

// parseHeader 读取消息头
func (c *Client) parseHeader() (err error) {
	var (
		n   int
		buf [512]byte
	)

	n, err = c.conn.Read(buf[:])
	if err != nil && err != io.EOF {
		log.Println(err.Error())
		return
	}

	var errMsg []byte
	if !c.Buffer.Validation(buf, n, &errMsg) {
		err = errors.New(string(errMsg))
		c.handleError(packet.NewErrMsg(errMsg))
		return
	}

	copy(c.Buffer.Header[:], buf[:n])

	return
}

// Body 读入消息体
func (c *Client) parseBody() (err error) {
	var (
		n   int
		buf [512]byte
	)
	length := c.Buffer.Length()
	for {
		n, err = c.conn.Read(buf[:])
		if err != nil && err != io.EOF {
			log.Println(err.Error())
		}

		c.Buffer.Body = append(c.Buffer.Body, buf[:n]...)
		length -= n

		if length < 0 {
			log.Println("报文头错误")
			c.handleError(packet.NewErrMsg([]byte("报文头错误")))
			return
		}

		if length == 0 {
			log.Println("Receive Data From Server: ", string(c.Buffer.Body[:]), len(string(c.Buffer.Body[:])))
			return
		}
	}
}

// handleError 错误处理函数
func (c *Client) handleError(errMsg *packet.ErrMsg) {
	c.Send(&packet.Msg{ErrMsg: errMsg})
}

// Dial connects to the remote server. Returns error if it can't.
func (c *Client) Dial(service string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		return err
	}
	c.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	return err
}

// Close close the socket
func (c *Client) Close() error {
	return c.conn.Close()
}

// Read dispatches incoming messages
func (c *Client) Read() {
	err := c.parseHeader()
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = c.parseBody()
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Receive message: ", string(c.Buffer.Body[:]))
}

// Send Send msg to server
func (c *Client) Send(msg *packet.Msg) {
	c.Buffer = packet.RPacket()

	buf, err := msg.Bytes()
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.Buffer.Body = append(c.Buffer.Body, buf[:]...)
	c.Buffer.SetLength(uint64(len(c.Buffer.Body)))
	log.Println(c.Buffer.Header)
	n, err := c.conn.Write(append(c.Buffer.Header[:], buf[:]...))
	if err != nil {
		log.Println("Write: ", err.Error())
		return
	}

	var msg2 packet.Msg
	var buffer bytes.Buffer
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&msg); err != nil {
		if err != io.EOF {
			log.Println(err.Error())
			return
		}
	}
	log.Printf("Send %d string: %+v\n", n, msg2)
}

// Run Run
func (c *Client) Run(service string) {
	// show line numbers
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if err := c.Dial(service); err != nil {
		log.Fatalln(err.Error())
	}

	defer func() {
		if err := c.Close(); err != err {
			log.Println(err.Error())
		}
	}()

	c.Send(&packet.Msg{Message: "sanghai"})

	// 	for {
	// 		if err = c.Read(); err != nil {
	// 			return err
	// 		}
	// 	}
}

// DefaultClient return a websocket hander
func DefaultClient() *Client {
	return &Client{
		// conn:,
		Buffer: packet.RPacket(),
		// Error:,
	}
}
