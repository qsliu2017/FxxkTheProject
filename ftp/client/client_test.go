package client

import (
	"archive/tar"
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8964"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")
			}
		}()
	}

	client, _ := NewFtpClient("localhost:8964")
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestLogin(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8965"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
					server.Writer.PrintfLine("331 User name okay, need password.")
				}

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
					server.Writer.PrintfLine("230 User logged in, proceed.")
				}
			}
		}()
	}

	client, _ := NewFtpClient("localhost:8965")
	if client == nil {
		t.Fatal("client is nil")
	}

	if err := client.Login("user", "pass"); err != nil {
		t.Fatal(err)
	}
}

func TestUsernameNotExist(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8966"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
					server.Writer.PrintfLine("332 Need account for login.")
				}
			}
		}()
	}

	client, _ := NewFtpClient("localhost:8966")
	if client == nil {
		t.Fatal("client is nil")
	}

	if err := client.Login("user", "pass"); err == nil || !errors.Is(err, ErrUsernameNotExist) {
		t.Fatal("should not login")
	}
}

func TestPasswordError(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8967"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
					server.Writer.PrintfLine("331 User name okay, need password.")
				}

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
					server.Writer.PrintfLine("530 Not logged in.")
				}
			}
		}()
	}

	client, _ := NewFtpClient("localhost:8967")
	if client == nil {
		t.Fatal("client is nil")
	}

	if err := client.Login("user", "pass"); err == nil || !errors.Is(err, ErrPasswordNotMatch) {
		t.Fatal("should not login")
	}
}

func TestLogout(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8968"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
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
			}
		}()
	}

	client, err := NewFtpClient("localhost:8968")
	if err != nil {
		t.Fatal(err)
	}
	if err = client.Login("user", "pass"); err != nil {
		t.Fatal(err)
	}
	if err := client.Logout(); err != nil {
		t.Fatal(err)
	}
}

func TestMode(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8969"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
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
			}
		}()
	}

	client, err := NewFtpClient("localhost:8969")
	if err != nil {
		t.Fatal(err)
	}
	if err = client.Mode(ModeStream); err != nil ||
		client.(*clientImpl).mode != ModeStream {
		t.Fatal(err)
	}
	if err = client.Mode(ModeBlock); err == nil ||
		!errors.Is(err, ErrModeNotSupported) ||
		client.(*clientImpl).mode != ModeStream {
		t.Fatal("should not change mode")
	}
	if err = client.Mode(ModeCompressed); err != nil ||
		client.(*clientImpl).mode != ModeCompressed {
		t.Fatal(err)
	}
}

func TestRetrStreamMode(t *testing.T) {
	os.Mkdir("_test_", 0777)
	defer os.RemoveAll("_test_")

	if listener, err := net.Listen("tcp", ":8970"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
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

					} else if strings.HasPrefix(line, "RETR large.txt") {
					} else if strings.HasPrefix(line, "RETR dir") {
					} else if strings.HasPrefix(line, "RETR") {
						// file not exist
					}
				}
			}
		}()
	}

	client, err := NewFtpClient("localhost:8970")
	if err != nil {
		t.Fatal(err)
	}
	if err = client.Retrieve("_test_/small.txt", "small.txt"); err != nil {
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
	if listener, err := net.Listen("tcp", ":8971"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
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
			}
		}()
	}

	client, err := NewFtpClient("localhost:8971")
	if err != nil {
		t.Fatal(err)
	}
	if err = client.Store("test_files/small.txt", "small.txt"); err != nil {
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

func TestStoreMultiFilesStreamMode(t *testing.T) {
	os.Mkdir("_test_", 0777)
	defer os.RemoveAll("_test_")

	filesave := make(chan bool, 10)
	if listener, err := net.Listen("tcp", ":8972"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
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
					} else if strings.HasPrefix(line, "STOR") {
						if dataConn == nil {
							server.Writer.PrintfLine("150 File status okay; about to open data connection.")
							continue
						}
						server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
						f, _ := os.Create("_test_/" + strings.TrimPrefix(line, "STOR "))
						io.Copy(f, dataConn)
						f.Close()
						filesave <- true
						dataConn.Close()
					}
				}
			}
		}()

	}
	client, err := NewFtpClient("localhost:8972")
	if err != nil {
		t.Fatal(err)
	}
	if err = client.Store("test_files/", ""); err != nil {
		t.Fatal(err)
	}

	remoteFiles, _ := ioutil.ReadDir("_test_")
	localFiles, _ := ioutil.ReadDir("test_files")
	for i := 0; i < len(localFiles); i++ {
		<-filesave
	}
	if len(remoteFiles) != len(localFiles) {
		t.Fatalf("remote files not equal to local files")
	}

	for i := range remoteFiles {
		local, _ := os.OpenFile("test_files/"+localFiles[i].Name(), os.O_RDONLY, 0666)
		defer local.Close()
		remote, _ := os.OpenFile("_test_/"+remoteFiles[i].Name(), os.O_RDONLY, 0666)
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
}

func TestStoreMultiFilesCompressedMode(t *testing.T) {
	os.Mkdir("_test_", 0777)
	defer os.RemoveAll("_test_")

	filesave := make(chan bool)
	if listener, err := net.Listen("tcp", ":8973"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
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
					} else if strings.HasPrefix(line, "STOR") {
						if dataConn == nil {
							server.Writer.PrintfLine("150 File status okay; about to open data connection.")
							continue
						}
						server.Writer.PrintfLine("125 Data connection already open; transfer starting.")
						f, _ := os.Create("_test_/" + strings.TrimPrefix(line, "STOR "))
						io.Copy(f, dataConn)
						f.Close()
						filesave <- true
						dataConn.Close()
					} else if strings.HasPrefix(line, "MODE") {
						server.Writer.PrintfLine("200 Command okay.")
					}
				}
			}
		}()
	}

	client, err := NewFtpClient("localhost:8973")
	if err != nil {
		t.Fatal(err)
	}

	client.Mode(ModeCompressed)
	if err = client.Store("test_files", "test_files.tar"); err != nil {
		t.Fatal(err)
	}

	<-filesave

	tarF, _ := os.Open("_test_/test_files.tar")
	defer tarF.Close()
	hasher := md5.New()
	io.Copy(hasher, tarF)
	remoteMd5 := hasher.Sum(nil)
	hasher.Reset()

	localFs, _ := ioutil.ReadDir("test_files")
	tarW := tar.NewWriter(hasher)
	for _, localF := range localFs {
		local, _ := os.OpenFile("test_files/"+localF.Name(), os.O_RDONLY, 0666)
		hdr, _ := tar.FileInfoHeader(localF, localF.Name())
		tarW.WriteHeader(hdr)
		io.Copy(tarW, local)
		local.Close()
	}
	tarW.Flush()
	tarW.Close()
	localMd5 := hasher.Sum(nil)

	if !bytes.Equal(localMd5, remoteMd5) {
		t.Fatalf("file not equal \n%x\n%x", localMd5, remoteMd5)
	}

}
