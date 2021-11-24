package client

import (
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
	ctrlConn, err := createCtrlConn(addr)
	if err != nil {
		return nil, err
	}

	client := defaultFtpClient()
	client.ctrlConn = ctrlConn

	return client, nil
}

func defaultFtpClient() *clientImpl {
	return &clientImpl{
		ctrlConn: nil,
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
	username string
	connMode byte
	mode     byte
	type_    byte
	stru     byte
}
