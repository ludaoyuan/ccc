package tcp

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

func handleTCPClient(conn net.Conn) {
	defer conn.Close()

	rPacket, wPacket := packet.PacketsCreation()

		bodyErr := rPacket.Validation(buf, n, &wPacket.Body)
		if bodyErr != nil {
			log.Println(wPacket.Body)
			conn.Write(append(wPacket.Header[:], wPacket.Body[:]...))
			conn.Close()
			return
		}

		length := rPacket.GetLength()

		log.Println(buf, "------------------------------------------")

		for {
			n, err = conn.Read(buf[:])
			if err != nil {
				// log.Println(err.Error())
				if err != io.EOF {
					log.Println(err.Error())
					conn.Close()
					return
				}

				log.Println("EOF")
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

				rPacket.ResetRPacket()

				wPacket.Body = []byte(`{"success":true,"data":"反馈信息"}`)
				bm := make(packet.Body)
				bm["code"] = "0000"
				bm["msg"] = "Success"

				var buf bytes.Buffer
				gob.Register(packet.Body{})
				encoder := gob.NewEncoder(&buf)
				if err := encoder.Encode(bm); err != nil {
					log.Println(err.Error())
					wPacket.ResetWPacket()
					break
				}
				bytes, err := ioutil.ReadAll(&buf)
				if err != nil {
					log.Println(err.Error())
					break
				}
				wPacket.Body = append(wPacket.Body, bytes[:]...)
				wPacket.SetLength()

				conn.Write(append(wPacket.Header[:], wPacket.Body[:]...))
				wPacket.ResetWPacket()
				break
			}
		}
	}
}

func Run(domain, port string) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
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
