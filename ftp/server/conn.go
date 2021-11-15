package server

import (
	"fmt"
	"ftp/cmd"
	"io"
	"log"
	"net"
	"strings"
)

type _FtpConn struct {
	ctrl     net.Conn
	data     net.Conn
	username string
	login    bool
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Println("Close connect", conn.RemoteAddr().String())
	}()
	ftpConn := _FtpConn{
		ctrl: conn,
		data: nil,
	}
	ftpConn.reply(cmd.SERVICE_READY, "Service ready for new user.")
	buf := make([]byte, 128)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Println(err)
			}
		}
		commandline := string(buf[:n])
		command := commandline[:4]
		if handler, has := commandHandlers[command]; has {
			if ftpConn._TestSyntax(commandline, handler.ArgsPattern, handler.Args...) {
				handler.Handler(&ftpConn, handler.Args...)
			}
		} else {
			ftpConn.reply(cmd.SYNTAX_ERROR, "Syntax error, command unrecognized.")
		}
	}
}

func (conn _FtpConn) reply(code int, msg string) error {
	if strings.Contains(msg, "\r\n") {
		return fmt.Errorf("multiline msg not implement")
	}
	_, err := conn.ctrl.Write([]byte(fmt.Sprintf("%3d %s\r\n", code, msg)))
	return err
}

func (conn _FtpConn) _SyntaxError() error {
	return conn.reply(cmd.SYNTAX_ERROR_IN_PARAM, "Syntax error in parameters or arguments.")
}

func (conn _FtpConn) _TestSyntax(cmd, syntax string, val ...interface{}) bool {
	_, err := fmt.Sscanf(cmd, syntax, val...)
	if err != nil {
		conn._SyntaxError()
		return false
	}
	return true
}
