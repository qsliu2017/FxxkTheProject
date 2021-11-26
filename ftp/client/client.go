package client

import (
	"net"
	"net/textproto"
)

type FtpClient interface {
	Login(username, password string) error
	Logout() error
	GetUsername() string

	ConnMode(byte) error
	GetConnMode() byte

	Mode(mode byte) error
	GetMode() byte

	Type(type_ byte) error
	GetType() byte

	Structure(stru byte) error
	GetStructure() byte

	Store(local, remote string) error
	Retrieve(local, remote string) error
}

func NewFtpClient(addr string) (FtpClient, error) {
	client := defaultFtpClient()

	if err := client.createCtrlConn(addr); err != nil {
		return nil, err
	}

	return client, nil
}

func defaultFtpClient() *clientImpl {
	return &clientImpl{
		ctrlConn: nil,
		dataConn: nil,
		username: "",
		connMode: ConnPort,
		mode:     ModeStream,
		type_:    TypeAscii,
		stru:     StruFile,
	}
}

var _ FtpClient = (*clientImpl)(nil)

type clientImpl struct {
	ctrlConn *textproto.Conn
	dataConn net.Conn
	username string
	connMode byte
	mode     byte
	type_    byte
	stru     byte
}

func (client *clientImpl) cmd(expect int, cmd string, args ...interface{}) (int, string, error) {
	if _, err := client.ctrlConn.Cmd(cmd, args...); err != nil {
		return 0, "", err
	}

	return client.ctrlConn.ReadResponse(expect)
}
