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
	if err := client.ctrlConn.Writer.PrintfLine("USER %s", username); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.USERNAME_OK); err != nil {
		switch code {
		case cmd.NEED_ACCOUNT:
			return ErrUsernameNotExist
		}
		return err
	}

	if err := client.ctrlConn.Writer.PrintfLine("PASS %s", password); err != nil {
		return err
	}

	if code, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.LOGIN_PROCEED); err != nil {
		switch code {
		case cmd.NOT_LOGIN:
			return ErrPasswordNotMatch
		}
		return err
	}

	client.username = username

	return nil
}

func (client *clientImpl) Logout() error {
	if err := client.ctrlConn.Writer.PrintfLine("QUIT"); err != nil {
		return err
	}

	if _, _, err := client.ctrlConn.Reader.ReadCodeLine(cmd.CTRL_CONN_CLOSE); err != nil {
		return err
	}

	client.username = ""

	return nil
}

func (client clientImpl) GetUsername() string {
	return client.username
}
