package client

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"io/ioutil"
	"log"
	"net"
	"packet"
)

var (
	conn    net.Conn
	rPacket *packet.Packet
	wPacket *packet.Packet
)

func checkFalt(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func handleSend(conn net.Conn, pkt *packet.Packet) {
	var buf bytes.Buffer

	gob.Register(packet.Body{})

	bm := make(packet.Body)
	encoder := gob.NewEncoder(&buf)

	bm["type"] = 0
	bm["data"] = "sanghai"

	err := encoder.Encode(bm)
	if err != nil {
		log.Println(err.Error())
		return
	}

	body, err := ioutil.ReadAll(&buf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	pkt.Body = append(pkt.Body, body[:]...)
	pkt.SetLength()
	_, err = conn.Write(append(pkt.Header[:], pkt.Body[:]...))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(pkt.Header)
}

func handle(conn net.Conn, pkt *packet.Packet, code uint, msg string) error {
	defer pkt.ResetRPacket()

	var buf bytes.Buffer
	bm := make(packet.Body)

	bm["type"] = 1
	bm["code"] = code
	bm["msg"] = msg

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(bm)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	body, err := ioutil.ReadAll(&buf)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(msg)
	pkt.Body = append(pkt.Body, body[:]...)
	_, err = conn.Write(append(pkt.Header[:], pkt.Body[:]...))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func handleRead(conn net.Conn, rPacket *packet.Packet) {
	var (
		buf [512]byte
		err error
		n   int
	)

	for {
		n, err = conn.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
				conn.Close()
			}
			break
		}

		bodyErr := rPacket.Validation(buf, n, &wPacket.Body)
		if bodyErr != nil {
			log.Println(string(wPacket.Body))
			conn.Write(append(wPacket.Header[:], wPacket.Body[:]...))
			conn.Close()
			return
		}

		length := int(binary.BigEndian.Uint64(buf[8:16]))

		log.Println(buf, "------------------------------------------")
		for {
			n, err = conn.Read(buf[:])
			if err != nil {
				log.Println(err.Error())
				if err != io.EOF {
					log.Println(err.Error())
					conn.Close()
					return
				}

				rPacket.ResetRPacket()
				break
			}
			rPacket.Body = append(rPacket.Body, buf[:n]...)

			length -= n

			if length < 0 {
				log.Println("报文头错误")
				conn.Close()
				return
			}

			if length == 0 {
				log.Println(string(rPacket.Body))
				rPacket.ResetRPacket()

				// err := handle(conn, wPacket, 0, "Success 反馈信息")
				// if err != nil {
				// 	log.Println(err.Error())
				// }
				break
			}
		}
	}
}

func connect(service string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkFalt(err)

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	checkFalt(err)
}

func Run(service string) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	connect(service)
	rPacket, wPacket = packet.PacketsCreation()

	handleSend(conn, wPacket)
	handleRead(conn, rPacket)
}
