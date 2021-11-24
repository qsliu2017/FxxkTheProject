package client

import (
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"testing"
)

func TestConnMode(t *testing.T) {
	serverConnChan := make(chan byte, 2)
	go func() {
		listener, _ := net.Listen("tcp", ":8969")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASV") {
				serverConnChan <- 0 // Got a PASV command from client
				dataConnListener, _ := net.Listen("tcp", ":0")
				addr := dataConnListener.Addr().(*net.TCPAddr)
				ip, port := []byte(addr.IP), addr.Port
				server.Writer.PrintfLine("227 Entering Passive Mode (%d,%d,%d,%d,%d,%d)", ip[0], ip[1], ip[2], ip[3], port>>8, port&0xff)
				dataConn, _ := dataConnListener.Accept()
				dataConnListener.Close()
				serverConnChan <- 1 //Got a data connection from PASV port
				dataConn.Close()
			} else if strings.HasPrefix(line, "PORT") {
				serverConnChan <- 2 // Got a PORT command from client
				var h1, h2, h3, h4, p1, p2 int
				fmt.Sscanf(line, "PORT %d,%d,%d,%d,%d,%d", &h1, &h2, &h3, &h4, &p1, &p2)
				dataConn, _ := net.Dial("tcp", fmt.Sprintf("%d.%d.%d.%d:%d", h1, h2, h3, h4, p1<<8+p2))
				server.Writer.PrintfLine("200 Command okay.")
				serverConnChan <- 3 // Got a data connection from PORT port
				dataConn.Close()
			}
		}
	}()

	c, _ := NewFtpClient("localhost:8969")
	client := c.(*clientImpl)

	client.createDataConn()
	if <-serverConnChan != 2 {
		t.Fatal("should get PORT command from client")
	}
	if <-serverConnChan != 3 {
		t.Fatal("should get data connection from PORT command")
	}

	c.ConnMode(ConnPasv)
	client.createDataConn()
	if <-serverConnChan != 0 {
		t.Fatal("should get PASV command from client")
	}
	if <-serverConnChan != 1 {
		t.Fatal("should get data connection from PASV port")
	}
}
