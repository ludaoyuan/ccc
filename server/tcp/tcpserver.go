package tcp

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"packet"
)

func handleTCPClient(conn net.Conn) {
	defer conn.Close()

	rPacket := &packet.Packet{[512]byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55}, []byte{}}
	copy(rPacket.Header[504:], []byte{0xff, 0xff, 0x55, 0x55, 0xff, 0xff, 0x55, 0x55})
	wPacket := &packet.Packet{[512]byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff}, []byte{}}
	copy(wPacket.Header[504:], []byte{0x55, 0x55, 0xff, 0xff, 0x55, 0x55, 0xff, 0xff})

	rPacket.Body = make([]byte, 0)

	var buf [512]byte
	var err error
	n := 0

	for {
		n, err = conn.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
				conn.Close()
			}
			break
		}

		if !rPacket.Validation(buf, n, &wPacket.Body) {
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
				if err != io.EOF {
					log.Println(err.Error())
					conn.Close()
					return
				}

				rPacket.Body = []byte{}
				break
			}

			rPacket.Body = append(rPacket.Body, buf[:n]...)

			length -= n

			if length < 0 {
				log.Println("报文头错误")
				conn.Write([]byte("报文头错误"))
				rPacket.Body = []byte{}
				conn.Close()
				return
			}

			if length == 0 {
				log.Println(string(rPacket.Body))

				rPacket.Body = []byte{}

				wPacket.Body = []byte(`{"success":true,"data":"反馈信息"}`)
				binary.BigEndian.PutUint64(wPacket.Header[8:16], uint64(len(rPacket.Body)))

				conn.Write(append(rPacket.Header[:], rPacket.Body[:]...))
				break
			}
		}
	}
}

func Run(domain, port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":9000")
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(tcpAddr)

	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println(conn.LocalAddr(), conn.RemoteAddr())

		go handleTCPClient(conn)
	}
}
