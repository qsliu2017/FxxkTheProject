package cmd

const (
	RESTART                   = 110
	NOT_READY                 = 120
	ALREADY_OPEN              = 125
	ABOUT_TO_DATA_CONN        = 150
	OK                        = 200
	_                         = 202
	_                         = 211
	_                         = 212
	_                         = 213
	_                         = 214
	_                         = 215
	SERVICE_READY             = 220
	CTRL_CONN_CLOSE           = 221
	_                         = 225
	_                         = 226
	_                         = 227
	LOGIN_PROCEED             = 230
	StatusFileActionCompleted = 250
	_                         = 257
	USERNAME_OK               = 331
	NEED_ACCOUNT              = 332
	_                         = 350
	NOT_AVAILABLE             = 421
	_                         = 425
	_                         = 426
	_                         = 450
	_                         = 451
	_                         = 452
	SYNTAX_ERROR              = 500
	SYNTAX_ERROR_IN_PARAM     = 501
	_                         = 502
	BAD_SEQUENCE              = 503
	_                         = 504
	NOT_LOGIN                 = 530
	NEED_ACCOUNT_FOR_STOR     = 532
	_                         = 550
	_                         = 551
	_                         = 552
	_                         = 553
)

var codeMessages = map[int]string{
	ALREADY_OPEN:              "Data connection already open; transfer starting.",
	ABOUT_TO_DATA_CONN:        "File status okay; about to open data connection.",
	StatusFileActionCompleted: "Requested file action okay, completed.",
}

func GetCodeMessage(code int) string {
	return codeMessages[code]
}
