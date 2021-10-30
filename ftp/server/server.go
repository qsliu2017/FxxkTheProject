package server

import (
	"fmt"
	"ftp/cmd"
	"log"
	"net"
	"strings"
)

type FtpConn struct {
	ctrl     net.Conn
	data     net.Conn
	username string
	login    bool
}

func Listen(port int) {
	ctrlConn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln(err)
	}
	defer ctrlConn.Close()

	conns := make(chan net.Conn)
	go func() {
		for {
			conn, err := ctrlConn.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("Accept connect from", conn.RemoteAddr().String(),
				"to", conn.LocalAddr().String())
			conns <- conn
		}
	}()

	for {
		go handleConn(<-conns)
	}

}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Println("Close connect", conn.RemoteAddr().String())
	}()
	ftp := FtpConn{
		ctrl: conn,
		data: nil,
	}
	ftp.reply(cmd.SERVICE_READY, "Service ready for new user.")
	buf := make([]byte, 128)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		command := string(buf[:n])
		switch {
		case strings.HasPrefix(command, "QUIT"):
			if ftp._TestSyntax(command, cmd.QUIT) {
				ftp.reply(cmd.CTRL_CONN_CLOSE, "Service closing control connection.")
				return
			}
		case strings.HasPrefix(command, "NOOP"):
			if ftp._TestSyntax(command, cmd.NOOP) {
				ftp.reply(200, "Command okay.")
			}
		case strings.HasPrefix(command, "USER"):
			log.Println("Accept command USER")
			var username string
			if ftp._TestSyntax(command, cmd.USER, &username) {
				log.Println("Parse command USER with username", username)
				if ftp.login {
					ftp.reply(cmd.LOGIN_PROCEED, "User logged in, proceed")
				} else if hasUser(username) {
					ftp.username = username
					ftp.reply(cmd.USERNAME_OK, "User name okay, need password.")
				} else {
					ftp.username = ""
					ftp.reply(cmd.NEED_ACCOUNT, "Need account for login.")
				}
			}
		case strings.HasPrefix(command, "PASS"):
			if ftp.login {

			} else if ftp.username == "" {
				ftp.reply(cmd.BAD_SEQUENCE, "Bad sequence of commands.")
			} else {
				var password string
				if ftp._TestSyntax(command, cmd.PASS, &password) {
					if testUser(ftp.username, password) {
						ftp.login = true
						ftp.reply(cmd.LOGIN_PROCEED, "User logged in, proceed.")
					} else {
						ftp.login = false
						ftp.username = ""
						ftp.reply(cmd.NOT_LOGIN, "Not logged in.")
					}
				}
			}
		}
	}
}

func (conn FtpConn) reply(code int, msg string) error {
	if strings.Contains(msg, "\r\n") {
		return fmt.Errorf("multiline msg not implement")
	}
	_, err := conn.ctrl.Write([]byte(fmt.Sprintf("%3d %s\r\n", code, msg)))
	return err
}

func (conn FtpConn) _SyntaxError() error {
	return conn.reply(cmd.SYNTAX_ERROR, "Syntax error\r\n")
}

func (conn FtpConn) _TestSyntax(cmd, syntax string, val ...interface{}) bool {
	_, err := fmt.Sscanf(cmd, syntax, val...)
	if err != nil {
		conn._SyntaxError()
		return false
	}
	return true
}
