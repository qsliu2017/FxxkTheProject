package server

import (
	"net"
)

type FtpServer interface {
	Listen(port int) error
	Close() error
	SetRootDir(string)
}

func NewFtpServer() FtpServer {
	return &_ServerImpl{}
}

var _ FtpServer = (*_ServerImpl)(nil)

type _ServerImpl struct {
	laddr    *net.TCPAddr
	listener *net.TCPListener
	rootDir  string
	// handlers map[chan<- bool]struct{} // notify all handlers to stop
}

func (server *_ServerImpl) Listen(port int) error {
	server.laddr = &net.TCPAddr{
		Port: port,
	}
	if listener, err := net.ListenTCP("tcp", server.laddr); err != nil {
		logger.Printf("failed to listen on %s: %s", server.laddr, err)
		return err
	} else {
		logger.Printf("server start listening on %s", server.laddr)
		server.listener = listener
		// server.handlers = make(map[chan<- bool]struct{})
		go func() {
			for {
				if conn, err := server.listener.Accept(); err != nil {
					if err == net.ErrClosed {
						return
					}
					// logger.Println(err)
				} else {
					logger.Printf("accepted connection from %s", conn.RemoteAddr())
					// channel := make(chan bool)
					// server.handlers[channel] = struct{}{}
					go handleClient(conn, server.rootDir) //Todo: pass channel to goroutine
				}
			}
		}()
		return nil
	}
}

func (server *_ServerImpl) Close() error {
	if err := server.listener.Close(); err != nil {
		logger.Printf("failed to close listener: %s", err)
		return err
	} else {
		logger.Printf("server listener closed")
		server.listener = nil
		server.laddr = nil
		// for channel := range server.handlers {
		// 	channel <- true
		// }
		// server.handlers = nil
		return nil
	}
}

func (server *_ServerImpl) SetRootDir(dir string) {
	server.rootDir = dir
}
