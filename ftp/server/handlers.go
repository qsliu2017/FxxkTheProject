package server

import (
	"ftp/cmd"
	"io"
	"net"
	"os"
)

type _RequestHandler func(conn *_FtpConn, args ...interface{}) error

type _CommandHandler struct {
	Handler     _RequestHandler
	ArgsPattern string
	Args        []interface{}
}

var commandHandlers map[string]_CommandHandler

func init() {
	commandHandlers = make(map[string]_CommandHandler)

	commandHandlers["QUIT"] = _CommandHandler{
		Handler:     quitHandler,
		ArgsPattern: cmd.QUIT,
		Args:        []interface{}{},
	}

	commandHandlers["NOOP"] = _CommandHandler{
		Handler:     noopHandler,
		ArgsPattern: cmd.NOOP,
		Args:        []interface{}{},
	}

	var user_username string
	commandHandlers["USER"] = _CommandHandler{
		Handler:     userHandler,
		ArgsPattern: cmd.USER,
		Args:        []interface{}{&user_username},
	}

	var pass_password string
	commandHandlers["PASS"] = _CommandHandler{
		Handler:     passHandler,
		ArgsPattern: cmd.PASS,
		Args:        []interface{}{&pass_password},
	}

	var port_h1, port_h2, port_h3, port_h4, port_p1, port_p2 byte
	commandHandlers["PORT"] = _CommandHandler{
		Handler:     portHandler,
		ArgsPattern: cmd.PORT,
		Args:        []interface{}{&port_h1, &port_h2, &port_h3, &port_h4, &port_p1, &port_p2},
	}

	var mode_modecode byte
	commandHandlers["MODE"] = _CommandHandler{
		Handler:     modeHandler,
		ArgsPattern: cmd.MODE,
		Args:        []interface{}{&mode_modecode},
	}

	var stor_pathname string
	commandHandlers["STOR"] = _CommandHandler{
		Handler:     storHandler,
		ArgsPattern: cmd.STOR,
		Args:        []interface{}{&stor_pathname},
	}

	var retr_pathname string
	commandHandlers["RETR"] = _CommandHandler{
		Handler:     retrHandler,
		ArgsPattern: cmd.RETR,
		Args:        []interface{}{&retr_pathname},
	}
}

var quitHandler _RequestHandler = func(conn *_FtpConn, args ...interface{}) error {
	return conn.reply(cmd.CTRL_CONN_CLOSE, "Service closing control connection.")
}

var noopHandler _RequestHandler = func(conn *_FtpConn, args ...interface{}) error {
	return conn.reply(cmd.OK, "Command okay.")
}

var userHandler _RequestHandler = func(conn *_FtpConn, args ...interface{}) error {
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

var passHandler _RequestHandler = func(conn *_FtpConn, args ...interface{}) error {
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

var portHandler _RequestHandler = func(conn *_FtpConn, args ...interface{}) error {
	// if conn.login {
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
	// } else {
	// 	return conn.reply(cmd.NOT_LOGIN, "Not logged in.")
	// }
}

var modeHandler _RequestHandler = func(conn *_FtpConn, args ...interface{}) error {
	mode := args[0].(*byte)
	switch *mode {
	case ModeStream:
		conn.mode = ModeStream
		return conn.reply(cmd.OK, "Command okay.")
	case ModeBlock:
		return conn.reply(cmd.StatusParamNotImplemented, cmd.GetCodeMessage(cmd.StatusParamNotImplemented))
	case ModeCompressed:
		conn.mode = ModeCompressed
		return conn.reply(cmd.OK, "Command okay.")
	default:
		return conn.reply(cmd.SYNTAX_ERROR_IN_PARAM, "Syntax error in parameters or arguments.")
	}
}

var storHandler _RequestHandler = func(conn *_FtpConn, args ...interface{}) error {
	if !conn.login {
		return conn.reply(cmd.NEED_ACCOUNT_FOR_STOR, "Need account for storing files.")
	}

	if conn.data == nil {
		return conn.reply(cmd.ABOUT_TO_DATA_CONN, cmd.GetCodeMessage(cmd.ABOUT_TO_DATA_CONN))
	}

	f, err := os.OpenFile(*args[0].(*string), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		//TODO: handle error
		return err
	}
	defer f.Close()
	conn.reply(cmd.ALREADY_OPEN, cmd.GetCodeMessage(cmd.ALREADY_OPEN))

	if _, err := io.Copy(f, conn.data); err != nil {
		return err
	}

	return conn.reply(cmd.StatusFileActionCompleted, cmd.GetCodeMessage(cmd.StatusFileActionCompleted))
}

var retrHandler _RequestHandler = func(conn *_FtpConn, args ...interface{}) error {
	if conn.data == nil {
		return conn.reply(cmd.ABOUT_TO_DATA_CONN, cmd.GetCodeMessage(cmd.ABOUT_TO_DATA_CONN))
	}

	f, err := os.Open(*args[0].(*string))
	if err != nil {
		//TODO
		return err
	}
	defer f.Close()
	conn.reply(cmd.ALREADY_OPEN, cmd.GetCodeMessage(cmd.ALREADY_OPEN))

	if _, err := io.Copy(conn.data, f); err != nil {
		return err
	}

	conn.data.Close()

	return conn.reply(cmd.StatusFileActionCompleted, cmd.GetCodeMessage(cmd.StatusFileActionCompleted))
}
