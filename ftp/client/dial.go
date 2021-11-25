package client

import (
	"ftp/cmd"
	"net"
	"net/textproto"
	"regexp"
	"strconv"
)

func (client *clientImpl) createCtrlConn(addr string) error {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return err
	}

	if _, _, err = conn.Reader.ReadCodeLine(cmd.SERVICE_READY); err != nil {
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

	addr := dataConnListener.Addr().(*net.TCPAddr)
	ip, port := []byte(addr.IP.To4()), addr.Port
	if err := client.ctrlConn.Writer.PrintfLine(
		cmd.PORT,
		ip[0], ip[1], ip[2], ip[3],
		(port >> 8), (port & 0xff)); err != nil {
		return nil, err
	}

	dataConn, err := dataConnListener.Accept()
	if err != nil {
		return nil, err
	}

	if code, _, err := client.ctrlConn.ReadCodeLine(200); err != nil {
		switch code {
		}
		return nil, err
	}

	return dataConn, nil
}

func (client *clientImpl) pasvDataConn() (net.Conn, error) {
	if err := client.ctrlConn.Writer.PrintfLine(cmd.PASV); err != nil {
		return nil, err
	}

	code, msg, err := client.ctrlConn.Reader.ReadCodeLine(cmd.StatusEnteringPasvMode)
	if err != nil {
		switch code {
		}
		return nil, err
	}

	addr, err := parsePasvResponse(msg)
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp4", nil, &addr)
}

var AddrRegexp = regexp.MustCompile(`\(([0-9]+,[0-9]+,[0-9]+,[0-9]+),([0-9]+),([0-9]+)\)`)

func parsePasvResponse(msg string) (net.TCPAddr, error) {
	matches := AddrRegexp.FindStringSubmatch(msg)
	if len(matches) != 4 {
		return net.TCPAddr{}, ErrInvalidPasvResponse
	}
	ip := net.ParseIP(matches[1])
	high, _ := strconv.Atoi(matches[2])
	low, _ := strconv.Atoi(matches[3])
	port := (high << 8) | low
	return net.TCPAddr{IP: ip, Port: port}, nil
}
