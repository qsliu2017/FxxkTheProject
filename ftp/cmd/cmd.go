package cmd

const (
	// USER<SP><username><CRLF>
	//  <username>::=<string>
	USER = "USER %s\r\n"

	// PASS<SP><password><CRLF>
	//  <password>::=<string>
	PASS = "PASS %s\r\n"

	// PASV<CRLF>
	PASV = "PASV\r\n"

	// QUIT<CRLF>
	QUIT = "QUIT\r\n"

	// PORT<SP><host-port><CRLF>
	//  <host-port>::=<host-number>,<port-number>
	//  <host-number>::=<number>,<number>,<number>,<number>
	//  <port-number>::=<number>,<number>
	PORT = "PORT %d,%d,%d,%d,%d,%d\r\n"

	// TYPE<SP><type-code><CRLF>
	//  <type-code>::=A[<SP><form-code>]
	//               |E[<SP><form-code>]
	//               |I
	//               |L<SP><byte-size>]
	//  <form-code>::=N|T|C
	TYPE = "TYPE %s\r\n"

	// MODE<SP><mode-code><CRLF>
	//  <mode-code>::=S|B|C
	MODE = "MODE %d\r\n"

	// STRU<SP><structure-code><CRLF>
	//  <structure-code>::=F|R|P
	STRU = "STRU %c\r\n"

	//RETR<SP><pathname><CRLF>
	RETR = "RETR %s\r\n"

	//STOR<SP><pathname><CRLF>
	STOR = "STOR %s\r\n"

	//NOOP<CRLF>
	NOOP = "NOOP\r\n"
)
