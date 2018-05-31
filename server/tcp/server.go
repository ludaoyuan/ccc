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

	rPacket, wPacket := packet.PacketsCreation()
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

		ok, bm := rPacket.Validation(buf, n, &wPacket.Body)
		if !ok {
			wPacket.SetLength()
			conn.Write(append(wPacket.Header[:], wPacket.Body[:]...))
			conn.Close()
			return
		}

		length := rPacket.GetLength()

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

				bm["code"] = uint16(0x0002)
				bm["msg"] = "Wrong header"

				wPacket.SetLength()
				conn.Write(append(wPacket.Header[:], wPacket.Body[:]...))
				conn.Close()

				rPacket.Body = []byte{}

				return
			}

			if length == 0 {
				log.Println(string(rPacket.Body))

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
