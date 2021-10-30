package server

import (
	"fmt"
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
func setupConn(t *testing.T) net.Conn {
	c, s := net.Pipe()
	go handleConn(s)

	// After connection establishment, expects 220
	assertReply(t, c, "220 Service ready for new user.\r\n", "Service not ready")
	return c
}

func endupConn(t *testing.T, c net.Conn) {
	c.Write([]byte(cmd.QUIT))
	assertReply(t, c, "221 Service closing control connection.\r\n", "Service quit error")

	c.Close()
}

func Test_Quit(t *testing.T) {
	c := setupConn(t)
	defer c.Close()

	c.Write([]byte(cmd.QUIT))
	assertReply(t, c, "221 Service closing control connection.\r\n", "Service quit error")
}

func Test_Noop(t *testing.T) {
	c := setupConn(t)
	defer endupConn(t, c)

	c.Write([]byte(cmd.NOOP))
	assertReply(t, c, "200 Command okay.\r\n", "Noop error")
}

func Test_User_with_valid_username(t *testing.T) {
	c := setupConn(t)
	defer endupConn(t, c)

	c.Write([]byte(fmt.Sprintf(cmd.USER, "test")))
	assertReply(t, c, "331 User name okay, need password.\r\n", "test valid user name error")
}

func Test_User_with_invalid_username(t *testing.T) {
	c := setupConn(t)
	defer endupConn(t, c)

	c.Write([]byte(fmt.Sprintf(cmd.USER, "test1")))
	assertReply(t, c, "332 Need account for login.\r\n", "test invalid user name error")
}

func Test_Pass_without_username(t *testing.T) {
	c := setupConn(t)
	defer endupConn(t, c)

	c.Write([]byte(fmt.Sprintf(cmd.PASS, "test")))
	assertReply(t, c, "503 Bad sequence of commands.\r\n", "test pass bad sequence error")
}

func Test_Pass_with_valid_account(t *testing.T) {
	c := setupConn(t)
	defer endupConn(t, c)

	c.Write([]byte(fmt.Sprintf(cmd.USER, "test")))
	assertReply(t, c, "331 User name okay, need password.\r\n", "test valid user name error")

	c.Write([]byte(fmt.Sprintf(cmd.PASS, "test")))
	assertReply(t, c, "230 User logged in, proceed.\r\n", "test valid account error")
}

func Test_Pass_with_invalid_account(t *testing.T) {
	c := setupConn(t)
	defer endupConn(t, c)

	c.Write([]byte(fmt.Sprintf(cmd.USER, "pikachu")))
	assertReply(t, c, "331 User name okay, need password.\r\n", "test valid user name error")

	c.Write([]byte(fmt.Sprintf(cmd.PASS, "pikachu")))
	assertReply(t, c, "530 Not logged in.\r\n", "test invalid account error")

	// After a failed login, username should be forgotten
	c.Write([]byte(fmt.Sprintf(cmd.PASS, "winnie")))
	assertReply(t, c, "503 Bad sequence of commands.\r\n", "test pass bad sequence error")
}
