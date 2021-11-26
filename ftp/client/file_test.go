package client

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"ftp/fm/block"
	"io"
	"net"
	"net/textproto"
	"os"
	"path"
	"strings"
	"testing"
)

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
			} else if strings.HasPrefix(line, "RETR ") {
				if dataConn == nil {
					server.Writer.PrintfLine("150 File status okay; about to open data connection.")
					continue
				}
				server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
				f, _ := os.Open(path.Join("test_files", line[len("RETR "):]))
				io.Copy(dataConn, f)
				dataConn.Close()
				f.Close()
				server.Writer.PrintfLine("250 Requested file action okay, completed.")
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8970")
	if err := client.Retrieve("_test_/small9993", "small9993"); err != nil {
		t.Fatal(err)
	}

	local, _ := os.Open("_test_/small9993")
	defer local.Close()
	remote, _ := os.Open("test_files/small9993")
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
			} else if strings.HasPrefix(line, "STOR ") {
				if dataConn == nil {
					server.Writer.PrintfLine("150 File status okay; about to open data connection.")
					continue
				}
				server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
				f, _ := os.Create(path.Join("_test_", line[len("STOR "):]))
				io.Copy(f, dataConn)
				server.Writer.PrintfLine("250 Requested file action okay, completed.")
				f.Close()
				dataConn.Close()
				filesave <- true

			}
		}
	}()

	client, _ := NewFtpClient("localhost:8971")
	if err := client.Store("test_files/small9993", "small9993"); err != nil {
		t.Fatal(err)
	}

	<-filesave
	// Don't go too fast

	local, _ := os.OpenFile("test_files/small9993", os.O_RDONLY, 0666)
	defer local.Close()
	remote, _ := os.OpenFile("_test_/small9993", os.O_RDONLY, 0666)
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

func TestStorBlockMode(t *testing.T) {
	os.Mkdir("_test_", 0777)
	defer os.RemoveAll("_test_")

	filesave := make(chan bool)
	go func() {
		listener, _ := net.Listen("tcp", ":8972")
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
			} else if strings.HasPrefix(line, "STOR ") {
				if dataConn == nil {
					server.Writer.PrintfLine("150 File status okay; about to open data connection.")
					continue
				}
				server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
				f, _ := os.Create(path.Join("_test_", line[len("STOR "):]))
				block.Receive(f, dataConn)
				server.Writer.PrintfLine("250 Requested file action okay, completed.")
				f.Close()
				filesave <- true
			} else if strings.HasPrefix(line, "MODE") {
				server.Writer.PrintfLine("200 Command okay.")
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8972")
	client.Mode(ModeBlock)
	if err := client.Store("test_files/small9993", "small9993"); err != nil {
		t.Fatal(err)
	}

	<-filesave
	// Don't go too fast

	local, _ := os.OpenFile("test_files/small9993", os.O_RDONLY, 0666)
	defer local.Close()
	remote, _ := os.OpenFile("_test_/small9993", os.O_RDONLY, 0666)
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

func TestRetrBlockMode(t *testing.T) {
	os.Mkdir("_test_", 0777)
	defer os.RemoveAll("_test_")

	filesave := make(chan bool)
	go func() {
		listener, _ := net.Listen("tcp", ":8973")
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
			} else if strings.HasPrefix(line, "RETR ") {
				if dataConn == nil {
					server.Writer.PrintfLine("150 File status okay; about to open data connection.")
					continue
				}
				server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
				f, _ := os.Open(path.Join("test_files", line[len("RETR "):]))
				block.Send(dataConn, f, 1<<8)
				server.Writer.PrintfLine("250 Requested file action okay, completed.")
				f.Close()
				filesave <- true
			} else if strings.HasPrefix(line, "MODE") {
				server.Writer.PrintfLine("200 Command okay.")
			}
		}
	}()

	client, _ := NewFtpClient("localhost:8973")
	client.Mode(ModeBlock)
	if err := client.Retrieve("_test_/small9993", "small9993"); err != nil {
		t.Fatal(err)
	}

	<-filesave
	// Don't go too fast

	local, _ := os.OpenFile("_test_/small9993", os.O_RDONLY, 0666)
	defer local.Close()
	remote, _ := os.OpenFile("test_files/small9993", os.O_RDONLY, 0666)
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
