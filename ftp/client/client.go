package client

type FtpClient interface {
	Login(username, password string) (string, error)
	Logout() (string, error)
	Mode(mode int) (string, error)
	Store(local, remote string) (string, error)
	Retrieve(local, remote string) (string, error)
}

func NewFtpClient(addr string) FtpClient {
	return nil
}

var _ FtpClient = (*clientImpl)(nil)

type clientImpl struct{}

func (*clientImpl) Login(username, password string) (string, error) {
	return "", nil
}

func (*clientImpl) Logout() (string, error) {
	return "", nil
}

const (
	ModeStream = iota
	ModeCompressed
)

func (*clientImpl) Mode(mode int) (string, error) {
	return "", nil
}

func (*clientImpl) Store(local, remote string) (string, error) {
	return "", nil
}

func (*clientImpl) Retrieve(local, remote string) (string, error) {
	return "", nil
}
