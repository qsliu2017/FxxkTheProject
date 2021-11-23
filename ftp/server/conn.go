package server

import (
	"fmt"
	"ftp/cmd"
	"io"
	"net"
	"net/textproto"
	"strings"
)

type _FtpConn struct {
	ctrl     *textproto.Conn
	data     net.Conn
	username string
	login    bool
	mode     byte
}

const (
	ModeStream     byte = 'S'
	ModeBlock      byte = 'B'
	ModeCompressed byte = 'C'
)

func handleConn(conn *net.Conn) {
	ftpConn := _FtpConn{
		ctrl: textproto.NewConn(*conn),
		data: nil,
		mode: ModeStream,
	}
	defer ftpConn.ctrl.Close()
	defer logger.Printf("connect %s closed\n", (*conn).RemoteAddr())

	ftpConn.reply(cmd.SERVICE_READY, "Service ready for new user.")
	for {
		commandline, err := ftpConn.ctrl.ReadLine()

		if err == io.EOF {
			logger.Printf("encount a EOF")
			return
		}

		if err == io.ErrClosedPipe {
			logger.Printf("encount a ErrClosedPipe")
			return
		}

		if err != nil {
			// logger.Println(err)
			continue
		}

		command := commandline[:4]
		handler, has := commandHandlers[command]
		if !has {
			ftpConn.reply(cmd.SYNTAX_ERROR, "Syntax error, command unrecognized.")
			continue
		}

		if ftpConn._TestSyntax(commandline, handler.ArgsPattern, handler.Args...) {
			logger.Printf("accept command %s from %s", command, (*conn).RemoteAddr())
			handler.Handler(&ftpConn, handler.Args...)
		}
	}
}

func (conn *_FtpConn) reply(code int, msg string) error {
	if strings.Contains(msg, "\r\n") {
		return fmt.Errorf("multiline msg not implement")
	}
	return conn.ctrl.PrintfLine("%3d %s", code, msg)
}

func (conn *_FtpConn) _SyntaxError() error {
	return conn.reply(cmd.SYNTAX_ERROR_IN_PARAM, "Syntax error in parameters or arguments.")
}

func (conn *_FtpConn) _TestSyntax(cmd, syntax string, val ...interface{}) bool {
	_, err := fmt.Sscanf(cmd, syntax, val...)
	if err != nil {
		conn._SyntaxError()
		return false
	}
	return true
}
