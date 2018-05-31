package client

import (
	"net"
	"packet"
	"testing"
)

const service = "127.0.0.1:9000"

func Test_handleSend(t *testing.T) {
	connect(service)
	rPacket, wPacket = packet.PacketsCreation()
	type args struct {
		conn net.Conn
		pkt  *packet.Packet
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"00001", args{conn, rPacket}},
		{"00002", args{conn, rPacket}},
		{"00003", args{conn, rPacket}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleSend(tt.args.conn, tt.args.pkt)
		})
	}
}

// func Test_handleRead(t *testing.T) {
// 	connect(service)
// 	rPacket, wPacket = packet.PacketsCreation()
// 	type args struct {
// 		conn    net.Conn
// 		rPacket *packet.Packet
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 		{"00001", args{conn, rPacket}},
// 		{"00002", args{conn, rPacket}},
// 		{"00003", args{conn, rPacket}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			handleRead(tt.args.conn, tt.args.rPacket)
// 		})
// 	}
// }
