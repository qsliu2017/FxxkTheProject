package server

import "fmt"

const (
	StatusTransferStarted = 125
	StatusFileStatusOK    = 150

	StatusOK                  = 200
	StatusReady               = 220
	StatusCloseConn           = 221
	StatusEnteringPasv        = 227
	StatuLoginProceed         = 230
	StatusFileActionCompleted = 250

	StatusUsernameOKNeedPassword = 331
	StatusNeedAccountForLogin    = 332

	StatusSyntaxError                        = 500
	StatusSyntaxErrorInParametersOrArguments = 501
	StatusCommandNotImplementedForParameter  = 504
	StatusNotLoggedIn                        = 530
	StatusFileUnavailable                    = 550
	StatusRequestedFileActionAborted         = 551
)

var ErrUnknownCode = fmt.Errorf("unknown code")

var codeMessages = map[int]string{
	StatusTransferStarted: "Data connection already open; transfer starting.",
	StatusFileStatusOK:    "File status okay; about to open data connection.",

	StatusOK:                  "Command okay.",
	StatusReady:               "Service ready for new user.",
	StatusCloseConn:           "Service closing control connection.",
	StatusEnteringPasv:        "Entering Passive Mode (%s).",
	StatuLoginProceed:         "User logged in, proceed.",
	StatusFileActionCompleted: "Requested file action okay, completed.",

	StatusUsernameOKNeedPassword: "User name okay, need password.",
	StatusNeedAccountForLogin:    "Need account for login.",

	StatusSyntaxError:                        "Syntax error, command unrecognized.",
	StatusSyntaxErrorInParametersOrArguments: "Syntax error in parameters or arguments.",
	StatusCommandNotImplementedForParameter:  "Command not implemented for that parameter.",
	StatusNotLoggedIn:                        "Not logged in.",
	StatusFileUnavailable:                    "File unavailable.",
	StatusRequestedFileActionAborted:         "Requested file action aborted, file unavailable.",
}

func (c *clientHandler) reply(code int, args ...interface{}) error {
	if msg, has := codeMessages[code]; has {
		resp := fmt.Sprintf("%d %s", code, fmt.Sprintf(msg, args...))
		logger.Printf("reply %s %s", c.username, resp)
		return c.ctrl.PrintfLine(resp)
	}
	return ErrUnknownCode
}
