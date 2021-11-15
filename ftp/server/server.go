package server

import (
	"fmt"
	"log"
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
		return err.Error(), err
	} else {
		server.listener = listener
		// server.handlers = make(map[chan<- bool]struct{})
		go func() {
			for {
				if conn, err := server.listener.Accept(); err != nil {
					if err == net.ErrClosed {
						return
					}
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

func (server *_ServerImpl) Close() (string, error) {
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
