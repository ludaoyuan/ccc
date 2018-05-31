package client

import (
	"packet"
	"testing"
)

func BenchmarkhandleSend(b *testing.B) {
	connect(service)
	rPacket, _ = packet.PacketsCreation()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handleSend(conn, rPacket)
	}
}
