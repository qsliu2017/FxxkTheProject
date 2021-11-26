package client

import (
	"ftp/cmd"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

func (client *clientImpl) createCtrlConn(addr string) error {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return err
	}

	if _, _, err = conn.Reader.ReadResponse(cmd.SERVICE_READY); err != nil {
		return err
	}

	client.ctrlConn = conn

	return nil
}

func (client *clientImpl) createDataConn() (err error) {
	if client.dataConn != nil {
		return nil
	}

	var conn net.Conn
	switch client.connMode {
	case ConnPasv:
		conn, err = client.pasvDataConn()
	case ConnPort:
		conn, err = client.portDataConn()
	default:
		err = ErrConnModeNotSupported
	}

	client.dataConn = conn

	return

}

func (client *clientImpl) closeDataConn() (err error) {
	if client.dataConn != nil {
		err = client.dataConn.Close()
		client.dataConn = nil
	}
	return
}

func (client *clientImpl) portDataConn() (net.Conn, error) {
	dataConnListener, err := net.ListenTCP("tcp4", nil)
	if err != nil {
		return nil, err
	}
	defer dataConnListener.Close()

	addr := dataConnListener.Addr().String()
	idx := strings.Index(addr, ":")
	hostParts := strings.Split(addr[:idx], ".")
	portNum, _ := strconv.Atoi(addr[idx+1:])

	if _, _, err := client.cmd(cmd.OK,
		"PORT %s,%d,%d",
		strings.Join(hostParts, ","),
		portNum>>8, portNum&0xff); err != nil {
		return nil, err
	}

	dataConn, err := dataConnListener.Accept()
	if err != nil {
		return nil, err
	}

	return dataConn, nil
}

func (client *clientImpl) pasvDataConn() (net.Conn, error) {
	_, msg, err := client.cmd(cmd.StatusEnteringPasvMode, cmd.PASV)
	if err != nil {
		return nil, err
	}

	start, end := strings.Index(msg, "("), strings.LastIndex(msg, ")")
	if start == -1 || end == -1 {
		return nil, ErrInvalidPasvResponse
	}

	data := strings.Split(msg[start+1:end], ",")
	if len(data) != 6 {
		return nil, ErrInvalidPasvResponse
	}

	host := strings.Join(data[:4], ".")

	portPart1, err := strconv.Atoi(data[4])
	if err != nil {
		return nil, err
	}

	portPart2, err := strconv.Atoi(data[5])
	if err != nil {
		return nil, err
	}

	port := (portPart1 << 8) | portPart2

	addr := net.JoinHostPort(host, strconv.Itoa(port))

	return net.Dial("tcp4", addr)
}
