package client

import (
	"net"
	"net/textproto"
	"testing"
)

func TestNewClient(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8964")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		server.Writer.PrintfLine("220 Service ready for new user.")
		server.Close()
	}()

	client, _ := NewFtpClient("localhost:8964")
	if client == nil {
		t.Fatal("client is nil")
	}
}
