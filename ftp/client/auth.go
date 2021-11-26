package client

import (
	"errors"
	"ftp/cmd"
)

var (
	ErrUsernameNotExist = errors.New("username does not exist")
	ErrPasswordNotMatch = errors.New("password does not match")
)

func (client *clientImpl) Login(username, password string) error {
	if code, msg, err := client.cmd(cmd.USERNAME_OK, "USER %s", username); err != nil {
		if code == cmd.NEED_ACCOUNT {
			return ErrUsernameNotExist
		}
		return errors.New(msg)
	}

	if code, msg, err := client.cmd(cmd.LOGIN_PROCEED, "PASS %s", password); err != nil {
		if code == cmd.NOT_LOGIN {
			return ErrPasswordNotMatch
		}
		return errors.New(msg)
	}

	client.username = username

	return nil
}

func (client *clientImpl) Logout() error {
	if _, msg, err := client.cmd(cmd.CTRL_CONN_CLOSE, "QUIT"); err != nil {
		return errors.New(msg)
	}

	client.username = ""

	return nil
}

func (client clientImpl) GetUsername() string {
	return client.username
}
