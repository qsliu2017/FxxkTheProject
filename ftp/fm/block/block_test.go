package block

import (
	"net"
	"os"
	"path"
	"testing"
)

func TestBlock(t *testing.T) {
	send := make(chan bool, 3)
	t.Run("sender", func(t *testing.T) {
		t.Parallel()

		sender, _ := net.Dial("tcp", "localhost:8964")
		defer sender.Close()

		dirName := "test_file"
		dir, _ := os.ReadDir(dirName)
		for _, file := range dir {
			f, _ := os.Open(path.Join(dirName, file.Name()))
			send <- true
			err := Send(sender, f, 1<<10)
			if err != nil {
				t.Fatal(err)
			}
		}
		send <- false

	})

	t.Run("receiver", func(t *testing.T) {
		t.Parallel()

		listener, _ := net.Listen("tcp", ":8964")
		recevicer, _ := listener.Accept()
		listener.Close()
		defer recevicer.Close()

		for <-send {
			f, _ := os.CreateTemp("", "ftp_fm_block_test_receiver_*")
			defer os.Remove(f.Name())

			if err := Receive(f, recevicer); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func BenchmarkBlock(b *testing.B) {

	go func() {
		sender, _ := net.Dial("tcp", "localhost:8964")
		defer sender.Close()

		dirName := "test_file"
		dir, _ := os.ReadDir(dirName)
		for i := 0; i < b.N; i++ {
			f, _ := os.Open(path.Join(dirName, dir[i%len(dir)].Name()))

			b.StartTimer()
			Send(sender, f, 1<<10)
			b.StopTimer()
		}
	}()

	listener, _ := net.Listen("tcp", ":8964")
	recevicer, _ := listener.Accept()
	listener.Close()
	defer recevicer.Close()

	for i := 0; i < b.N; i++ {
		f, _ := os.CreateTemp("", "ftp_fm_block_test_receiver_*")
		defer os.Remove(f.Name())

		if err := Receive(f, recevicer); err != nil {
			b.Fatal(err)
		}
	}
}
