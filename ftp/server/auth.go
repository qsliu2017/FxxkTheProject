package server

import "errors"

var (
	ErrCloseConn                = errors.New("close connection")
	_            commandHandler = (*clientHandler).handleUSER
	_            commandHandler = (*clientHandler).handlePASS
	_            commandHandler = (*clientHandler).handleQUIT

	account = map[string]string{
		"test":      "test",
		"pikachu":   "winnie",
		"anonymous": "anonymous",
	}
)

func (c *clientHandler) handleUSER(param string) error {
	if c.login {
		return c.reply(StatuLoginProceed)
	}

	if _, has := account[param]; has {
		c.username = param
		return c.reply(StatusUsernameOKNeedPassword)
	} else {
		c.username = ""
		return c.reply(StatusNeedAccountForLogin)
	}
}

func (c *clientHandler) handlePASS(param string) error {
	if password := account[c.username]; password == param {
		c.login = true
		return c.reply(StatuLoginProceed)
	} else {
		c.login = false
		c.username = ""
		return c.reply(StatusNotLoggedIn)
	}
}

func (c *clientHandler) handleQUIT(param string) error {
	c.reply(StatusCloseConn)
	return ErrCloseConn
}
