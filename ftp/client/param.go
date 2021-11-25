package client

import (
	"errors"
	"ftp/cmd"
)

const (
	ConnPasv byte = iota
	ConnPort

	ModeStream     byte = 'S'
	ModeBlock      byte = 'B'
	ModeCompressed byte = 'C'

	TypeAscii  byte = 'A'
	TypeBinary byte = 'I'

	StruFile byte = 'F'
)

var (
	ErrConnModeNotSupported = errors.New("connection mode not supported")
	ErrInvalidPasvResponse  = errors.New("invalid pasv response")
	ErrModeNotSupported     = errors.New("mode not support")
	ErrTypeNotSupported     = errors.New("type not support")
	ErrStruNotSupported     = errors.New("stru not support")
)

func (client *clientImpl) ConnMode(mode byte) error {
	if mode != ConnPasv && mode != ConnPort {
		return ErrConnModeNotSupported
	}
	client.connMode = mode
	return nil
}

func (client clientImpl) GetConnMode() byte {
	return client.connMode
}

func (client *clientImpl) Mode(mode byte) error {
	if mode != ModeStream && mode != ModeBlock && mode != ModeCompressed {
		return ErrModeNotSupported
	}

	if code, msg, err := client.cmd(cmd.OK, "MODE %c", mode); err != nil {
		if code == cmd.StatusParamNotImplemented {
			return ErrModeNotSupported
		}
		return errors.New(msg)
	}

	client.mode = mode

	return nil
}

func (client clientImpl) GetMode() byte {
	return client.mode
}

func (client *clientImpl) Type(type_ byte) error {
	if type_ != TypeAscii && type_ != TypeBinary {
		return ErrTypeNotSupported
	}

	if code, msg, err := client.cmd(cmd.OK, "TYPE %c", type_); err != nil {
		if code == cmd.StatusParamNotImplemented {
			return ErrTypeNotSupported
		}
		return errors.New(msg)
	}

	client.type_ = type_

	return nil
}

func (cleint clientImpl) GetType() byte {
	return cleint.type_
}

func (client *clientImpl) Structure(stru byte) error {
	if stru != StruFile {
		return ErrStruNotSupported
	}

	if code, msg, err := client.cmd(cmd.OK, "STRU %c", stru); err != nil {
		if code == cmd.StatusParamNotImplemented {
			return ErrStruNotSupported
		}
		return errors.New(msg)
	}

	client.stru = stru
	return nil
}

func (client clientImpl) GetStructure() byte {
	return client.stru
}
