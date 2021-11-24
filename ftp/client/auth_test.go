package client

import (
	"errors"
	"net"
	"net/textproto"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8965")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
			server.Writer.PrintfLine("331 User name okay, need password.")
		}

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
			server.Writer.PrintfLine("230 User logged in, proceed.")
		}

	}()

	client, _ := NewFtpClient("localhost:8965")

	if err := client.Login("user", "pass"); err != nil {
		t.Fatal(err)
	}
	if client.GetUsername() != "user" {
		t.Fatal("username is not user")
	}
}

func TestUsernameNotExist(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8966")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
			server.Writer.PrintfLine("332 Need account for login.")
		}
	}()

	client, _ := NewFtpClient("localhost:8966")

	if err := client.Login("user", "pass"); err == nil || !errors.Is(err, ErrUsernameNotExist) {
		t.Fatal("should not login")
	}
	if client.GetUsername() != "" {
		t.Fatal("username should be empty")
	}
}

func TestPasswordError(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8967")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
			server.Writer.PrintfLine("331 User name okay, need password.")
		}

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
			server.Writer.PrintfLine("530 Not logged in.")
		}
	}()

	client, _ := NewFtpClient("localhost:8967")

	if err := client.Login("user", "pass"); err == nil || !errors.Is(err, ErrPasswordNotMatch) {
		t.Fatal("should not login")
	}
	if client.GetUsername() != "" {
		t.Fatal("username should be empty")
	}
}

func TestLogout(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8968")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
			server.Writer.PrintfLine("331 User name okay, need password.")
		}

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
			server.Writer.PrintfLine("230 User logged in, proceed.")
		}

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "QUIT") {
			server.Writer.PrintfLine("221 Service closing control connection.")
		}
	}()

	client, _ := NewFtpClient("localhost:8968")
	if err := client.Login("user", "pass"); err != nil {
		t.Fatal(err)
	}
	if client.GetUsername() != "user" {
		t.Fatal("username is not user")
	}
	if err := client.Logout(); err != nil {
		t.Fatal(err)
	}
	if client.GetUsername() != "" {
		t.Fatal("username should be empty")
	}
}
