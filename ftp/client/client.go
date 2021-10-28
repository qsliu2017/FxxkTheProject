package client

import (
	"bytes"
	"fmt"
	"ftp/cmd"
	"log"
	"net"
)

func CtrlConnect(addr string) {
	ctrlConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer ctrlConn.Close()
	ctrlConn.Write([]byte(fmt.Sprintf(cmd.NOOP)))
	log.Println("sending NOOP")
	printResp(ctrlConn)
	ctrlConn.Write([]byte(fmt.Sprintf(cmd.QUIT)))
	log.Println("sending QUIT")
	printResp(ctrlConn)
}

func printResp(conn net.Conn) {
	buf := make([]byte, 64)
	for {
		n, _ := conn.Read(buf)
		resp := buf[:n]
		log.Println(resp)
		print(string(buf))
		if bytes.HasSuffix(resp, []byte("\r\n")) {
			return
		}
	}
}
