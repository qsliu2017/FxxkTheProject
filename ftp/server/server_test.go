package server

import (
	"fmt"
	"ftp/cmd"
	"net"
	"os"
	"strings"
	"testing"
)

var (
	__buffer = make([]byte, 1024)
	__n      int
)

// new a FtpServer, listen on a port, test if it can be dial
func Test_Listen(t *testing.T) {
	t.Run("valid port", func(t *testing.T) {
		s := NewFtpServer()
		port := 8964
		if _, err := s.Listen(port); err != nil {
			t.Error("Listen error", err)
		}
		defer s.Close()

		if conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port)); err != nil {
			t.Error("Dial error", err)
		} else {
			conn.Close()
		}
	})
}

func readReply(c net.Conn) {
	__n, _ = c.Read(__buffer)
}

func assertReply(t *testing.T, c net.Conn, expect, msg string) {
	t.Helper()
	if readReply(c); strings.Compare(expect, string(__buffer[:__n])) != 0 {
		t.Error(msg,
			"\nExpect:", []byte(expect),
			"\nActual:", __buffer[:__n],
		)
	}
}

// Setup a mock connection, test if service ready, and return the client conn.
func setupConn(t *testing.T) net.Conn {
	t.Helper()
	c, s := net.Pipe()
	go handleConn(s)

	// After connection establishment, expects 220
	assertReply(t, c, "220 Service ready for new user.\r\n", "Service not ready")
	return c
}

func teardownConn(t *testing.T, c net.Conn) {
	t.Helper()
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
	defer teardownConn(t, c)

	c.Write([]byte(cmd.NOOP))
	assertReply(t, c, "200 Command okay.\r\n", "Noop error")
}

func Test_User(t *testing.T) {
	t.Run("valid username", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.USER, "test")))
		assertReply(t, c, "331 User name okay, need password.\r\n", "test valid user name error")
	})

	t.Run("invalid username", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.USER, "test1")))
		assertReply(t, c, "332 Need account for login.\r\n", "test invalid user name error")
	})
}

func Test_Pass(t *testing.T) {
	t.Run("without username", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.PASS, "test")))
		assertReply(t, c, "503 Bad sequence of commands.\r\n", "test pass bad sequence error")
	})

	t.Run("valid account", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.USER, "test")))
		assertReply(t, c, "331 User name okay, need password.\r\n", "test valid user name error")

		c.Write([]byte(fmt.Sprintf(cmd.PASS, "test")))
		assertReply(t, c, "230 User logged in, proceed.\r\n", "test valid account error")
	})

	t.Run("invalid account", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.USER, "pikachu")))
		assertReply(t, c, "331 User name okay, need password.\r\n", "test valid user name error")

		c.Write([]byte(fmt.Sprintf(cmd.PASS, "pikachu")))
		assertReply(t, c, "530 Not logged in.\r\n", "test invalid account error")

		// After a failed login, username should be forgotten
		c.Write([]byte(fmt.Sprintf(cmd.PASS, "winnie")))
		assertReply(t, c, "503 Bad sequence of commands.\r\n", "test pass bad sequence error")
	})
}

func Test_Port(t *testing.T) {
	t.Run("not login", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.PORT, 0, 0, 0, 0, 0, 0)))
		assertReply(t, c, "530 Not logged in.\r\n", "test port with login error")
	})

	t.Run("login", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.USER, "test")))
		assertReply(t, c, "331 User name okay, need password.\r\n", "test valid user name error")
		c.Write([]byte(fmt.Sprintf(cmd.PASS, "test")))
		assertReply(t, c, "230 User logged in, proceed.\r\n", "test valid account error")

		// Client listens on a port
		dataConn, err := net.Listen("tcp", ":5456")
		if err != nil {
			t.Fatal(err)
		}
		defer dataConn.Close()

		accept := make(chan struct{})
		go func() {
			conn, err := dataConn.Accept()
			if err != nil {
				t.Log(err)
			}
			defer conn.Close()
			accept <- struct{}{}
		}()

		// Client sends the port to server
		c.Write([]byte(fmt.Sprintf(cmd.PORT, 127, 0, 0, 1, 21, 80)))
		assertReply(t, c, "200 Command okay.\r\n", "test port error")
		if _, has := <-accept; !has {
			t.Error("data port not accept any connection")
		}
	})
}

func Test_Mode(t *testing.T) {
	t.Run("", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.MODE, 'S')))
		assertReply(t, c, "200 Command okay.\r\n", "test mode error")

		c.Write([]byte(fmt.Sprintf(cmd.MODE, 'C')))
		assertReply(t, c, "200 Command okay.\r\n", "test mode error")

		c.Write([]byte(fmt.Sprintf(cmd.MODE, 'B')))
		assertReply(t, c, "504 Command not implemented for that parameter.\r\n", "test mode error")

		c.Write([]byte(fmt.Sprintf(cmd.MODE, 'A')))
		assertReply(t, c, "501 Syntax error in parameters or arguments.\r\n", "test mode error")
	})
}

func Test_Stor(t *testing.T) {
	t.Run("not login", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.STOR, "test.txt")))
		assertReply(t, c, "532 Need account for storing files.\r\n", "test stor with login error")
	})

	t.Run("data connect", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.USER, "test")))
		assertReply(t, c, "331 User name okay, need password.\r\n", "test valid user name error")
		c.Write([]byte(fmt.Sprintf(cmd.PASS, "test")))
		assertReply(t, c, "230 User logged in, proceed.\r\n", "test valid account error")

		c.Write([]byte(fmt.Sprintf(cmd.STOR, "test.txt")))
		assertReply(t, c, "150 File status okay; about to open data connection.\r\n", "")

		dataConn, err := net.Listen("tcp", ":5456")
		if err != nil {
			t.Fatal(err)
		}
		defer dataConn.Close()

		accept := make(chan net.Conn)
		go func() {
			conn, err := dataConn.Accept()
			if err != nil {
				t.Log(err)
			}
			accept <- conn
		}()

		// Client sends the port to server
		c.Write([]byte(fmt.Sprintf(cmd.PORT, 127, 0, 0, 1, 21, 80)))
		assertReply(t, c, "200 Command okay.\r\n", "test port error")
		dataChan, has := <-accept
		if !has {
			t.Error("data port not accept any connection")
		}

		c.Write([]byte(fmt.Sprintf(cmd.STOR, "test.txt")))
		assertReply(t, c, "125 Data connection already open; transfer starting.\r\n", "")

		dataChan.Write([]byte("test data\r\n"))
		dataChan.Close()
		assertReply(t, c, "250 Requested file action okay, completed.\r\n", "")

		f, _ := os.Open("test.txt")
		defer func() {
			f.Close()
			os.Remove("test.txt")
		}()
		buffer := make([]byte, 32)
		n, _ := f.Read(buffer)
		if strings.Compare(string(buffer[:n]), "test data\r\n") != 0 {
			t.Error("data not match")
		}
	})

	// t.Run("compressed mode")
}

func Test_Retr(t *testing.T) {
	t.Run("data connect", func(t *testing.T) {
		c := setupConn(t)
		defer teardownConn(t, c)

		c.Write([]byte(fmt.Sprintf(cmd.RETR, "test_retr.txt")))
		assertReply(t, c, "150 File status okay; about to open data connection.\r\n", "")

		dataConn, err := net.Listen("tcp", ":5456")
		if err != nil {
			t.Fatal(err)
		}
		defer dataConn.Close()

		accept := make(chan net.Conn)
		go func() {
			conn, err := dataConn.Accept()
			if err != nil {
				t.Log(err)
			}
			accept <- conn
		}()

		// Client sends the port to server
		c.Write([]byte(fmt.Sprintf(cmd.PORT, 127, 0, 0, 1, 21, 80)))
		assertReply(t, c, "200 Command okay.\r\n", "test port error")
		dataChan, has := <-accept
		if !has {
			t.Error("data port not accept any connection")
		}
		defer dataChan.Close()

		c.Write([]byte(fmt.Sprintf(cmd.STOR, "test.txt")))
		assertReply(t, c, "125 Data connection already open; transfer starting.\r\n", "")

		buffer := make([]byte, 32)
		file, _ := os.Create("test_retr_back.txt")
		defer file.Close()
		for {
			n, err := dataChan.Read(buffer)
			if err != nil {
				assertReply(t, c, "250 Requested file action okay, completed.\r\n", "")
				break
			}
			file.Write(buffer[:n])
		}
	})
}
