package server

import (
	"fmt"
	"ftp/cmd"
	"log"
	"net"
	"strings"
)

type FtpServer interface {
	Listen(port int) (string, error)
	Close() (string, error)
}

func NewFtpServer() FtpServer {
	return &serverImpl{}
}

var _ FtpServer = (*serverImpl)(nil)

type serverImpl struct {
	laddr    *net.TCPAddr
	listener *net.TCPListener
	// handlers map[chan<- bool]struct{} // notify all handlers to stop
}

func (server *serverImpl) Listen(port int) (string, error) {
	server.laddr = &net.TCPAddr{
		Port: port,
	}
	if listener, err := net.ListenTCP("tcp", server.laddr); err != nil {
		return err.Error(), err
	} else {
		server.listener = listener
		// server.handlers = make(map[chan<- bool]struct{})
		go func() {
			for {
				if conn, err := server.listener.Accept(); err != nil {
					log.Println(err)
				} else {
					// channel := make(chan bool)
					// server.handlers[channel] = struct{}{}
					go handleConn(conn) //Todo: pass channel to goroutine
				}
			}
		}()
		return fmt.Sprintf("Ftp server listening on %s", listener.Addr().String()), nil
	}
}

func (server *serverImpl) Close() (string, error) {
	if err := server.listener.Close(); err != nil {
		return err.Error(), err
	} else {
		server.listener = nil
		server.laddr = nil
		// for channel := range server.handlers {
		// 	channel <- true
		// }
		// server.handlers = nil
		return "Ftp server closed", nil
	}
}

type ftpConn struct {
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
	ftp := ftpConn{
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
			if ftp._TestSyntax(commandline, handler.ArgsPattern, handler.Args...) {
				handler.Handler(&ftp, handler.Args...)
			}
		} else {
			ftp.reply(cmd.SYNTAX_ERROR, "Syntax error, command unrecognized.")
		}
	}
}

func (conn ftpConn) reply(code int, msg string) error {
	if strings.Contains(msg, "\r\n") {
		return fmt.Errorf("multiline msg not implement")
	}
	_, err := conn.ctrl.Write([]byte(fmt.Sprintf("%3d %s\r\n", code, msg)))
	return err
}

func (conn ftpConn) _SyntaxError() error {
	return conn.reply(cmd.SYNTAX_ERROR_IN_PARAM, "Syntax error in parameters or arguments.")
}

func (conn ftpConn) _TestSyntax(cmd, syntax string, val ...interface{}) bool {
	_, err := fmt.Sscanf(cmd, syntax, val...)
	if err != nil {
		conn._SyntaxError()
		return false
	}
	return true
}
