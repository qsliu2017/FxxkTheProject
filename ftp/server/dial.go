package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	ErrConnectToDataPort                = errors.New("connect to data port failed")
	_                    commandHandler = (*clientHandler).handlePORT
	_                    commandHandler = (*clientHandler).handlePASV
)

func (c *clientHandler) handlePORT(param string) error {
	parts := strings.Split(param, ",")
	if len(parts) != 6 {
		return c.reply(StatusSyntaxErrorInParametersOrArguments)
	}

	ip := strings.Join(parts[:4], ".")
	portPart1, err1 := strconv.Atoi(parts[4])
	portPart2, err2 := strconv.Atoi(parts[5])
	if err1 != nil || err2 != nil {
		return c.reply(StatusSyntaxErrorInParametersOrArguments)
	}
	port := (portPart1 << 8) | portPart2

	conn, err := net.Dial("tcp", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return ErrConnectToDataPort
	}

	c.conn = conn

	c.data = bufio.NewReadWriter(
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	)

	return c.reply(StatusOK)
}

func (c *clientHandler) handlePASV(param string) error {
	listener, err := net.ListenTCP("tcp4", nil)
	if err != nil {
		return err
	}
	defer listener.Close()

	host, port, _ := net.SplitHostPort(listener.Addr().String())

	ip := strings.Split(host, ".")
	portPart1, _ := strconv.Atoi(port)
	portPart2 := portPart1 & 0xff
	portPart1 >>= 8

	if err := c.reply(StatusEnteringPasv, fmt.Sprintf("%s,%d,%d", strings.Join(ip, ","), portPart1, portPart2)); err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	c.conn = conn

	c.data = bufio.NewReadWriter(
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	)

	return nil
}
