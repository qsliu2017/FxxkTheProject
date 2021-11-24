package client

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8964")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		server.Writer.PrintfLine("220 Service ready for new user.")
		server.Close()
	}()

	client, _ := NewFtpClient("localhost:8964")
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestLogin(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8965")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
			server.Writer.PrintfLine("331 User name okay, need password.")
		}

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
			server.Writer.PrintfLine("230 User logged in, proceed.")
		}

	}()

	client, _ := NewFtpClient("localhost:8965")

	if err := client.Login("user", "pass"); err != nil {
		t.Fatal(err)
	}
	if client.GetUsername() != "user" {
		t.Fatal("username is not user")
	}
}

func TestUsernameNotExist(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8966")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
			server.Writer.PrintfLine("332 Need account for login.")
		}
	}()

	client, _ := NewFtpClient("localhost:8966")

	if err := client.Login("user", "pass"); err == nil || !errors.Is(err, ErrUsernameNotExist) {
		t.Fatal("should not login")
	}
	if client.GetUsername() != "" {
		t.Fatal("username should be empty")
	}
}

func TestPasswordError(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8967")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
			server.Writer.PrintfLine("331 User name okay, need password.")
		}

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
			server.Writer.PrintfLine("530 Not logged in.")
		}
	}()

	client, _ := NewFtpClient("localhost:8967")

	if err := client.Login("user", "pass"); err == nil || !errors.Is(err, ErrPasswordNotMatch) {
		t.Fatal("should not login")
	}
	if client.GetUsername() != "" {
		t.Fatal("username should be empty")
	}
}

func TestLogout(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8968")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
			server.Writer.PrintfLine("331 User name okay, need password.")
		}

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
			server.Writer.PrintfLine("230 User logged in, proceed.")
		}

		if line, _ := server.ReadLine(); strings.HasPrefix(line, "QUIT") {
			server.Writer.PrintfLine("221 Service closing control connection.")
		}
	}()

	client, _ := NewFtpClient("localhost:8968")
	if err := client.Login("user", "pass"); err != nil {
		t.Fatal(err)
	}
	if client.GetUsername() != "user" {
		t.Fatal("username is not user")
	}
	if err := client.Logout(); err != nil {
		t.Fatal(err)
	}
	if client.GetUsername() != "" {
		t.Fatal("username should be empty")
	}
}

func TestConnMode(t *testing.T) {
	serverConnChan := make(chan byte, 2)
	go func() {
		listener, _ := net.Listen("tcp", ":8969")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASV") {
				serverConnChan <- 0 // Got a PASV command from client
				dataConnListener, _ := net.Listen("tcp", ":0")
				addr := dataConnListener.Addr().(*net.TCPAddr)
				ip, port := []byte(addr.IP), addr.Port
				server.Writer.PrintfLine("227 Entering Passive Mode (%d,%d,%d,%d,%d,%d)", ip[0], ip[1], ip[2], ip[3], port>>8, port&0xff)
				dataConn, _ := dataConnListener.Accept()
				dataConnListener.Close()
				serverConnChan <- 1 //Got a data connection from PASV port
				dataConn.Close()
			} else if strings.HasPrefix(line, "PORT") {
				serverConnChan <- 2 // Got a PORT command from client
				var h1, h2, h3, h4, p1, p2 int
				fmt.Sscanf(line, "PORT %d,%d,%d,%d,%d,%d", &h1, &h2, &h3, &h4, &p1, &p2)
				dataConn, _ := net.Dial("tcp", fmt.Sprintf("%d.%d.%d.%d:%d", h1, h2, h3, h4, p1<<8+p2))
				server.Writer.PrintfLine("200 Command okay.")
				serverConnChan <- 3 // Got a data connection from PORT port
				dataConn.Close()
			}
		}
	}()

	c, _ := NewFtpClient("localhost:8969")
	client := c.(*clientImpl)

	client.createDataConn()
	if <-serverConnChan != 2 {
		t.Fatal("should get PORT command from client")
	}
	if <-serverConnChan != 3 {
		t.Fatal("should get data connection from PORT command")
	}

	c.ConnMode(ConnPasv)
	client.createDataConn()
	if <-serverConnChan != 0 {
		t.Fatal("should get PASV command from client")
	}
	if <-serverConnChan != 1 {
		t.Fatal("should get data connection from PASV port")
	}
}

func TestMode(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8969")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "MODE S") ||
				strings.HasPrefix(line, "MODE C") {
				server.Writer.PrintfLine("200 Command okay.")
			} else if strings.HasPrefix(line, "MODE B") {
				server.Writer.PrintfLine("504 Command not implemented for that parameter.")
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8969")
	if err := client.Mode(ModeStream); err != nil ||
		client.GetMode() != ModeStream {
		t.Fatal(err)
	}
	if err := client.Mode(ModeBlock); err == nil ||
		!errors.Is(err, ErrModeNotSupported) ||
		client.GetMode() != ModeStream {
		t.Fatal("should not change mode")
	}
	if err := client.Mode(ModeCompressed); err != nil ||
		client.GetMode() != ModeCompressed {
		t.Fatal(err)
	}
}

func TestType(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8970")
		conn, _ := listener.Accept()
		listener.Close()

		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "TYPE A") ||
				strings.HasPrefix(line, "TYPE I") {
				server.Writer.PrintfLine("200 Command okay.")
			} else if strings.HasPrefix(line, "TYPE L") {
				server.Writer.PrintfLine("504 Command not implemented for that parameter.")
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8970")

	if err := client.Type(TypeAscii); err != nil ||
		client.GetType() != TypeAscii {
		t.Fatal(err)
	}

	if err := client.Type('L'); err == nil ||
		!errors.Is(err, ErrTypeNotSupported) ||
		client.GetType() != TypeAscii {
		t.Fatal("should not change type")
	}

	if err := client.Type(TypeBinary); err != nil ||
		client.GetType() != TypeBinary {
		t.Fatal(err)
	}

}

func TestStru(t *testing.T) {
	go func() {
		listener, _ := net.Listen("tcp", ":8971")
		conn, _ := listener.Accept()
		listener.Close()

		server := textproto.NewConn(conn)
		defer server.Close()
		server.Writer.PrintfLine("220 Service ready for new user.")

		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "STRU F") {
				server.Writer.PrintfLine("200 Command okay.")
			} else if strings.HasPrefix(line, "STRU") {
				server.Writer.PrintfLine("504 Command not implemented for that parameter.")
			}
		}
	}()
	client, _ := NewFtpClient("localhost:8971")
	if err := client.Structure(StruFile); err != nil ||
		client.GetStructure() != StruFile {
		t.Fatal(err)
	}
	if err := client.Structure('S'); err == nil ||
		!errors.Is(err, ErrStruNotSupported) ||
		client.GetStructure() != StruFile {
		t.Fatal("should not change stru")
	}
}

func TestRetrStreamMode(t *testing.T) {
	os.Mkdir("_test_", 0777)
	defer os.RemoveAll("_test_")

	go func() {
		listener, _ := net.Listen("tcp", ":8970")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()

		var dataConn net.Conn

		server.Writer.PrintfLine("220 Service ready for new user.")
		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "PORT") {
				var h1, h2, h3, h4, p1, p2 byte
				fmt.Sscanf(line, "PORT %d,%d,%d,%d,%d,%d", &h1, &h2, &h3, &h4, &p1, &p2)
				dataConn, _ = net.Dial("tcp", fmt.Sprintf("%d.%d.%d.%d:%d", h1, h2, h3, h4, int(p1)*256+int(p2)))
				server.Writer.PrintfLine("200 Command okay.")
			} else if strings.HasPrefix(line, "RETR small.txt") {
				if dataConn == nil {
					server.Writer.PrintfLine("150 File status okay; about to open data connection.")
					continue
				}
				server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
				f, _ := os.Open("test_files/small.txt")
				io.Copy(dataConn, f)
				dataConn.Close()
				f.Close()
				server.Writer.PrintfLine("250 Requested file action okay, completed.")

			} else if strings.HasPrefix(line, "RETR large.txt") {
			} else if strings.HasPrefix(line, "RETR dir") {
			} else if strings.HasPrefix(line, "RETR") {
				// file not exist
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8970")
	if err := client.Retrieve("_test_/small.txt", "small.txt"); err != nil {
		t.Fatal(err)
	}

	local, _ := os.Open("_test_/small.txt")
	defer local.Close()
	remote, _ := os.Open("test_files/small.txt")
	defer remote.Close()
	hasher := md5.New()
	io.Copy(hasher, local)
	localMd5 := hasher.Sum(nil)
	hasher.Reset()
	io.Copy(hasher, remote)
	remoteMd5 := hasher.Sum(nil)
	if !bytes.Equal(localMd5, remoteMd5) {
		t.Fatal("file not equal")
	}
}

func TestStorStreamMode(t *testing.T) {
	os.Mkdir("_test_", 0777)
	defer os.RemoveAll("_test_")

	filesave := make(chan bool)
	go func() {
		listener, _ := net.Listen("tcp", ":8971")
		conn, _ := listener.Accept()
		listener.Close()
		server := textproto.NewConn(conn)
		defer server.Close()

		var dataConn net.Conn

		server.Writer.PrintfLine("220 Service ready for new user.")
		for {
			if line, _ := server.ReadLine(); strings.HasPrefix(line, "PORT") {
				var h1, h2, h3, h4, p1, p2 byte
				fmt.Sscanf(line, "PORT %d,%d,%d,%d,%d,%d", &h1, &h2, &h3, &h4, &p1, &p2)
				dataConn, _ = net.Dial("tcp", fmt.Sprintf("%d.%d.%d.%d:%d", h1, h2, h3, h4, int(p1)*256+int(p2)))
				server.Writer.PrintfLine("200 Command okay.")
			} else if strings.HasPrefix(line, "STOR small.txt") {
				if dataConn == nil {
					server.Writer.PrintfLine("150 File status okay; about to open data connection.")
					continue
				}
				server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
				f, _ := os.Create("_test_/small.txt")
				io.Copy(f, dataConn)
				f.Close()
				filesave <- true
				dataConn.Close()
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8971")
	if err := client.Store("test_files/small.txt", "small.txt"); err != nil {
		t.Fatal(err)
	}

	<-filesave
	// Don't go too fast

	local, _ := os.OpenFile("test_files/small.txt", os.O_RDONLY, 0666)
	defer local.Close()
	remote, _ := os.OpenFile("_test_/small.txt", os.O_RDONLY, 0666)
	defer remote.Close()

	hasher := md5.New()
	io.Copy(hasher, local)
	localMd5 := hasher.Sum(nil)
	hasher.Reset()
	io.Copy(hasher, remote)
	remoteMd5 := hasher.Sum(nil)
	if !bytes.Equal(localMd5, remoteMd5) {
		t.Fatalf("file not equal \n%x\n%x", localMd5, remoteMd5)
	}
}

// func TestStoreMultiFilesStreamMode(t *testing.T) {
// 	os.Mkdir("_test_", 0777)
// 	defer os.RemoveAll("_test_")

// 	filesave := make(chan bool, 10)
// 	go func() {
// 		listener, _ := net.Listen("tcp", ":8972")
// 		conn, _ := listener.Accept()
// 		listener.Close()
// 		server := textproto.NewConn(conn)
// 		defer server.Close()

// 		var dataConn net.Conn

// 		server.Writer.PrintfLine("220 Service ready for new user.")
// 		for {
// 			if line, _ := server.ReadLine(); strings.HasPrefix(line, "PORT") {
// 				var h1, h2, h3, h4, p1, p2 byte
// 				fmt.Sscanf(line, "PORT %d,%d,%d,%d,%d,%d", &h1, &h2, &h3, &h4, &p1, &p2)
// 				dataConn, _ = net.Dial("tcp", fmt.Sprintf("%d.%d.%d.%d:%d", h1, h2, h3, h4, int(p1)*256+int(p2)))
// 				server.Writer.PrintfLine("200 Command okay.")
// 			} else if strings.HasPrefix(line, "STOR") {
// 				if dataConn == nil {
// 					server.Writer.PrintfLine("150 File status okay; about to open data connection.")
// 					continue
// 				}
// 				server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
// 				f, _ := os.Create("_test_/" + strings.TrimPrefix(line, "STOR "))
// 				io.Copy(f, dataConn)
// 				f.Close()
// 				filesave <- true
// 				dataConn.Close()
// 			}
// 		}
// 	}()
// 	client, _ := NewFtpClient("localhost:8972")
// 	if err := client.Store("test_files/", ""); err != nil {
// 		t.Fatal(err)
// 	}

// 	remoteFiles, _ := ioutil.ReadDir("_test_")
// 	localFiles, _ := ioutil.ReadDir("test_files")
// 	for i := 0; i < len(localFiles); i++ {
// 		<-filesave
// 	}
// 	if len(remoteFiles) != len(localFiles) {
// 		t.Fatalf("remote files not equal to local files")
// 	}

// 	for i := range remoteFiles {
// 		local, _ := os.OpenFile("test_files/"+localFiles[i].Name(), os.O_RDONLY, 0666)
// 		defer local.Close()
// 		remote, _ := os.OpenFile("_test_/"+remoteFiles[i].Name(), os.O_RDONLY, 0666)
// 		defer remote.Close()

// 		hasher := md5.New()
// 		io.Copy(hasher, local)
// 		localMd5 := hasher.Sum(nil)
// 		hasher.Reset()
// 		io.Copy(hasher, remote)
// 		remoteMd5 := hasher.Sum(nil)
// 		if !bytes.Equal(localMd5, remoteMd5) {
// 			t.Fatalf("file not equal \n%x\n%x", localMd5, remoteMd5)
// 		}
// 	}
// }

// func TestStoreMultiFilesCompressedMode(t *testing.T) {
// 	os.Mkdir("_test_", 0777)
// 	defer os.RemoveAll("_test_")

// 	filesave := make(chan bool)
// 	go func() {
// 		listener, _ := net.Listen("tcp", ":8973")
// 		conn, _ := listener.Accept()
// 		listener.Close()
// 		server := textproto.NewConn(conn)
// 		defer server.Close()

// 		var dataConn net.Conn

// 		server.Writer.PrintfLine("220 Service ready for new user.")
// 		for {
// 			if line, _ := server.ReadLine(); strings.HasPrefix(line, "PORT") {
// 				var h1, h2, h3, h4, p1, p2 byte
// 				fmt.Sscanf(line, "PORT %d,%d,%d,%d,%d,%d", &h1, &h2, &h3, &h4, &p1, &p2)
// 				dataConn, _ = net.Dial("tcp", fmt.Sprintf("%d.%d.%d.%d:%d", h1, h2, h3, h4, int(p1)*256+int(p2)))
// 				server.Writer.PrintfLine("200 Command okay.")
// 			} else if strings.HasPrefix(line, "STOR") {
// 				if dataConn == nil {
// 					server.Writer.PrintfLine("150 File status okay; about to open data connection.")
// 					continue
// 				}
// 				server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
// 				f, _ := os.Create("_test_/" + strings.TrimPrefix(line, "STOR "))
// 				io.Copy(f, dataConn)
// 				f.Close()
// 				filesave <- true
// 				dataConn.Close()
// 			} else if strings.HasPrefix(line, "MODE") {
// 				server.Writer.PrintfLine("200 Command okay.")
// 			}
// 		}
// 	}()

// 	client, _ := NewFtpClient("localhost:8973")

// 	client.Mode(ModeCompressed)
// 	if err := client.Store("test_files", "test_files.tar"); err != nil {
// 		t.Fatal(err)
// 	}

// 	<-filesave

// 	tarF, _ := os.Open("_test_/test_files.tar")
// 	defer tarF.Close()
// 	hasher := md5.New()
// 	io.Copy(hasher, tarF)
// 	remoteMd5 := hasher.Sum(nil)
// 	hasher.Reset()

// 	localFs, _ := ioutil.ReadDir("test_files")
// 	tarW := tar.NewWriter(hasher)
// 	for _, localF := range localFs {
// 		local, _ := os.OpenFile("test_files/"+localF.Name(), os.O_RDONLY, 0666)
// 		hdr, _ := tar.FileInfoHeader(localF, localF.Name())
// 		tarW.WriteHeader(hdr)
// 		io.Copy(tarW, local)
// 		local.Close()
// 	}
// 	tarW.Flush()
// 	tarW.Close()
// 	localMd5 := hasher.Sum(nil)

// 	if !bytes.Equal(localMd5, remoteMd5) {
// 		t.Fatalf("file not equal \n%x\n%x", localMd5, remoteMd5)
// 	}

// }
