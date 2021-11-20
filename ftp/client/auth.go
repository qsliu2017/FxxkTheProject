package client

import (
	"errors"
)

var (
	ErrUsernameNotExist error = errors.New("username does not exist")
	ErrPasswordNotMatch error = errors.New("password does not match")
)
