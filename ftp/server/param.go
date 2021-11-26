package server

var (
	_ commandHandler = (*clientHandler).handleMODE
	_ commandHandler = (*clientHandler).handleTYPE
	_ commandHandler = (*clientHandler).handleSTRU
)

const (
	ModeStream     byte = 'S'
	ModeBlock      byte = 'B'
	ModeCompressed byte = 'C'

	TypeAscii  byte = 'A'
	TypeBinary byte = 'I'

	StruFile byte = 'F'
)

func (c *clientHandler) handleMODE(param string) error {
	mode := param[0]
	switch mode {
	case ModeStream, ModeBlock:
		c.mode = mode
		return c.reply(StatusOK)
	default:
		return c.reply(StatusCommandNotImplementedForParameter)
	}
}

func (c *clientHandler) handleTYPE(param string) error {
	type_ := param[0]
	switch type_ {
	case TypeAscii, TypeBinary:
		c.type_ = type_
		return c.reply(StatusOK)
	default:
		return c.reply(StatusCommandNotImplementedForParameter)
	}
}

func (c *clientHandler) handleSTRU(param string) error {
	stru := param[0]
	switch stru {
	case StruFile:
		c.stru = stru
		return c.reply(StatusOK)
	default:
		return c.reply(StatusCommandNotImplementedForParameter)
	}
}
