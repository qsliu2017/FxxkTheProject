package server

import (
	"fmt"
	"net"
)

type FtpServer interface {
	Listen(port int) (string, error)
	Close() (string, error)
}

func NewFtpServer() FtpServer {
	return &_ServerImpl{}
}

var _ FtpServer = (*_ServerImpl)(nil)

type _ServerImpl struct {
	laddr    *net.TCPAddr
	listener *net.TCPListener
	// handlers map[chan<- bool]struct{} // notify all handlers to stop
}

func (server *_ServerImpl) Listen(port int) (string, error) {
	server.laddr = &net.TCPAddr{
		Port: port,
	}
	if listener, err := net.ListenTCP("tcp", server.laddr); err != nil {
		logger.Printf("Failed to listen on %s: %s", server.laddr, err)
		return err.Error(), err
	} else {
		logger.Printf("Listening on %s", server.laddr)
		server.listener = listener
		// server.handlers = make(map[chan<- bool]struct{})
		go func() {
			for {
				if conn, err := server.listener.Accept(); err != nil {
					if err == net.ErrClosed {
						logger.Printf("Ftp server closed")
						return
					}
					logger.Println(err)
				} else {
					logger.Printf("Ftp server accepted connection from %s", conn.RemoteAddr())
					// channel := make(chan bool)
					// server.handlers[channel] = struct{}{}
					go handleConn(conn) //Todo: pass channel to goroutine
				}
			}
		}()
		return fmt.Sprintf("Ftp server listening on %s", listener.Addr().String()), nil
	}
}

func (server *_ServerImpl) Close() (string, error) {
	if err := server.listener.Close(); err != nil {
		logger.Printf("Failed to close listener: %s", err)
		return err.Error(), err
	} else {
		logger.Printf("Listener closed")
		server.listener = nil
		server.laddr = nil
		// for channel := range server.handlers {
		// 	channel <- true
		// }
		// server.handlers = nil
		return "Ftp server closed", nil
	}
}
