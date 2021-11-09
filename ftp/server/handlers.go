package server

import (
	"ftp/cmd"
	"net"
)

type RequestHandler func(conn *FtpConn, args ...interface{}) error

type commandHandler struct {
	Handler     RequestHandler
	ArgsPattern string
	Args        []interface{}
}

var commandHandlers map[string]commandHandler

func init() {
	commandHandlers = make(map[string]commandHandler)

	commandHandlers["QUIT"] = commandHandler{
		Handler:     quitHandler,
		ArgsPattern: cmd.QUIT,
		Args:        []interface{}{},
	}

	commandHandlers["NOOP"] = commandHandler{
		Handler:     noopHandler,
		ArgsPattern: cmd.NOOP,
		Args:        []interface{}{},
	}

	var user_username string
	commandHandlers["USER"] = commandHandler{
		Handler:     userHandler,
		ArgsPattern: cmd.USER,
		Args:        []interface{}{&user_username},
	}

	var pass_password string
	commandHandlers["PASS"] = commandHandler{
		Handler:     passHandler,
		ArgsPattern: cmd.PASS,
		Args:        []interface{}{&pass_password},
	}

	var port_h1, port_h2, port_h3, port_h4, port_p1, port_p2 byte
	commandHandlers["PORT"] = commandHandler{
		Handler:     portHandler,
		ArgsPattern: cmd.PORT,
		Args:        []interface{}{&port_h1, &port_h2, &port_h3, &port_h4, &port_p1, &port_p2},
	}
}

var quitHandler RequestHandler = func(conn *FtpConn, args ...interface{}) error {
	return conn.reply(cmd.CTRL_CONN_CLOSE, "Service closing control connection.")
}

var noopHandler RequestHandler = func(conn *FtpConn, args ...interface{}) error {
	return conn.reply(cmd.OK, "Command okay.")
}

var userHandler RequestHandler = func(conn *FtpConn, args ...interface{}) error {
	if conn.login {
		return conn.reply(cmd.LOGIN_PROCEED, "User logged in, proceed")
	} else if username, ok := args[0].(*string); ok && hasUser(*username) {
		conn.username = *username
		return conn.reply(cmd.USERNAME_OK, "User name okay, need password.")
	} else {
		conn.username = ""
		return conn.reply(cmd.NEED_ACCOUNT, "Need account for login.")
	}
}

var passHandler RequestHandler = func(conn *FtpConn, args ...interface{}) error {
	if conn.login {
		// do something
		return nil
	} else if conn.username == "" {
		return conn.reply(cmd.BAD_SEQUENCE, "Bad sequence of commands.")
	} else {
		if password, ok := args[0].(*string); ok && testUser(conn.username, *password) {
			conn.login = true
			return conn.reply(cmd.LOGIN_PROCEED, "User logged in, proceed.")
		} else {
			conn.login = false
			conn.username = ""
			return conn.reply(cmd.NOT_LOGIN, "Not logged in.")
		}
	}
}

var portHandler RequestHandler = func(conn *FtpConn, args ...interface{}) error {
	if conn.login {
		var err error
		conn.data, err = net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   net.IPv4(*args[0].(*byte), *args[1].(*byte), *args[2].(*byte), *args[3].(*byte)),
			Port: int(*args[4].(*byte))*256 + int(*args[5].(*byte)),
		})
		if err != nil {
			return err
		} else {
			return conn.reply(cmd.OK, "Command okay.")
		}
	} else {
		return conn.reply(cmd.NOT_LOGIN, "Not logged in.")
	}
}
