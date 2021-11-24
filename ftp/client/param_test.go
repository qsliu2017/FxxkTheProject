package client

import (
	"errors"
	"net"
	"net/textproto"
	"strings"
	"testing"
)

func TestMode(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8969")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "MODE S") ||
				strings.HasPrefix(line, "MODE C") {
				server.Writer.PrintfLine("200 Command okay.")
			} else if strings.HasPrefix(line, "MODE B") {
				server.Writer.PrintfLine("504 Command not implemented for that parameter.")
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8969")
	if err := client.Mode(ModeStream); err != nil ||
		client.GetMode() != ModeStream {
		t.Fatal(err)
	}
	if err := client.Mode(ModeBlock); err == nil ||
		!errors.Is(err, ErrModeNotSupported) ||
		client.GetMode() != ModeStream {
		t.Fatal("should not change mode")
	}
	if err := client.Mode(ModeCompressed); err != nil ||
		client.GetMode() != ModeCompressed {
		t.Fatal(err)
	}
}

func TestType(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8970")
		conn, _ := listener.Accept()
		listener.Close()

		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "TYPE A") ||
				strings.HasPrefix(line, "TYPE I") {
				server.Writer.PrintfLine("200 Command okay.")
			} else if strings.HasPrefix(line, "TYPE L") {
				server.Writer.PrintfLine("504 Command not implemented for that parameter.")
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8970")

	if err := client.Type(TypeAscii); err != nil ||
		client.GetType() != TypeAscii {
		t.Fatal(err)
	}

	if err := client.Type('L'); err == nil ||
		!errors.Is(err, ErrTypeNotSupported) ||
		client.GetType() != TypeAscii {
		t.Fatal("should not change type")
	}

	if err := client.Type(TypeBinary); err != nil ||
		client.GetType() != TypeBinary {
		t.Fatal(err)
	}

}

func TestStru(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8971")
		conn, _ := listener.Accept()
		listener.Close()

		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "STRU F") {
				server.Writer.PrintfLine("200 Command okay.")
			} else if strings.HasPrefix(line, "STRU") {
				server.Writer.PrintfLine("504 Command not implemented for that parameter.")
			}
		}
	}()
	client, _ := NewFtpClient("localhost:8971")
	if err := client.Structure(StruFile); err != nil ||
		client.GetStructure() != StruFile {
		t.Fatal(err)
	}
	if err := client.Structure('S'); err == nil ||
		!errors.Is(err, ErrStruNotSupported) ||
		client.GetStructure() != StruFile {
		t.Fatal("should not change stru")
	}
}
