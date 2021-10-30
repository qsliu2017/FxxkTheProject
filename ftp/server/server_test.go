package server

import (
	"ftp/cmd"
	"net"
	"strings"
	"testing"
)

var (
	__buffer = make([]byte, 1024)
	__n      int
)

func readReply(c net.Conn) {
	__n, _ = c.Read(__buffer)
}

func assertReply(t *testing.T, c net.Conn, expect, msg string) {
	if readReply(c); strings.Compare(expect, string(__buffer[:__n])) != 0 {
		t.Error(msg,
			"\nExpect:", []byte(expect),
			"\nActual:", __buffer[:__n],
		)
	}
}

// Setup a mock connection, test if service ready, and return the client conn.
func setup_conn(t *testing.T) net.Conn {
	c, s := net.Pipe()
	go handleConn(s)

	// After connection establishment, expects 220
	assertReply(t, c, "220 Service ready for new user.\r\n", "Service not ready")
	return c
}

func endup_conn(t *testing.T, c net.Conn) {
	c.Write([]byte(cmd.QUIT))
	assertReply(t, c, "221 Service closing control connection.\r\n", "Service quit error")

	c.Close()
}

func Test_Quit(t *testing.T) {
	c := setup_conn(t)
	defer c.Close()

	c.Write([]byte(cmd.QUIT))
	assertReply(t, c, "221 Service closing control connection.\r\n", "Service quit error")
}

func Test_Noop(t *testing.T) {
	c := setup_conn(t)
	defer endup_conn(t, c)

	c.Write([]byte(cmd.NOOP))
	assertReply(t, c, "200 Command okay.\r\n", "Noop error")
}
