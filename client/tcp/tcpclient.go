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

	_ "net/http/pprof"

	_ "github.com/mkevac/debugcharts"
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

		if ok, _ := rPacket.Validation(buf, n, &rPacket.Body); !ok {
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
					return
				}

				rPacket.ResetRPacket()
				break
			}
			rPacket.Body = append(rPacket.Body, buf[:n]...)

			length -= n

			if length < 0 {
				log.Println("报文头错误")
				conn.Write([]byte("报文头错误"))
				rPacket.Body = []byte{}
				conn.Close()
				rPacket.ResetRPacket()
				break
			}

			if length == 0 {
				log.Println(string(rPacket.Body))

				rPacket.Body = []byte{}

				wPacket.Body = []byte(`{"success":true,"data":"反馈信息"}`)
				binary.BigEndian.PutUint64(wPacket.Header[8:16], uint64(len(rPacket.Body)))

				conn.Write(append(rPacket.Header[:], rPacket.Body[:]...))
				break
			}

			rPacket.ResetRPacket()
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

	handleSend(conn, rPacket)
	handleRead(conn, rPacket)
}
