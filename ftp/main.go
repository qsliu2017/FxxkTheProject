package main

import (
	"flag"
	"ftp/client"
	"ftp/server"
	"strconv"
)

var (
	mode string
	addr string
)

func init() {
	flag.StringVar(&mode, "mode", "server", "server|client")
	flag.StringVar(&addr, "addr", "", "listened port of server or connect port of client")
	flag.Parse()
}
func main() {
	switch mode {
	case "s":
		fallthrough
	case "server":
		port, _ := strconv.Atoi(addr)
		server.Listen(port)
	case "c":
		fallthrough
	case "client":
		client.CtrlConnect(addr)
	}
}
