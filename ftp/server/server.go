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
		commandline := string(buf[:n])
		command := commandline[:4]
		if handler, has := commandHandlers[command]; has {
			ftp._TestSyntax(commandline, handler.ArgsPattern, handler.Args...)
			handler.Handler(&ftp, handler.Args...)
			continue
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
