package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

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
	buf := make([]byte, 128)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		req := buf[:n]
		command := string(req)
		switch {
		case strings.HasPrefix(command, "QUIT"):
			conn.Write([]byte("200 Connect closed\r\n"))
			return
		case strings.HasPrefix(command, "NOOP"):
			conn.Write([]byte("200 Do nothing\r\n"))
		}
	}
}
