package block

import (
	"encoding/binary"
	"errors"
	"hash"
	"hash/fnv"
	"io"
)

var _BlockSize int64 = 1 << 10

func SetBlockSize(size int64) {
	_BlockSize = size
}

// In block mode, each file is divided into blocks, and each block has a fixed size.
// The conn is shared by all the files, the header is to spilt the file.
type _BlockHdr struct {
	BlockSize int64
}

type _BlockFtr struct {
	Length   int64
	Checksum uint64
	EOF      bool
}

var HashAlgorithm = func() hash.Hash64 { return fnv.New64() }
var ErrBrokenBlock = errors.New("block broken")

//Given a stream of data, pack it into one block, and write it to the writer w.
func Send(dst io.Writer, src io.Reader) (err error) {
	binary.Write(dst, binary.BigEndian, _BlockHdr{_BlockSize})

	r := io.TeeReader(src, dst)

	blockFtr := _BlockFtr{EOF: false}

	hasher := HashAlgorithm()

	for {
		hasher.Reset()
		if blockFtr.Length, err = io.CopyN(hasher, r, _BlockSize); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		blockFtr.Checksum = hasher.Sum64()

		binary.Write(dst, binary.BigEndian, blockFtr)
	}

	padding := make([]byte, _BlockSize-blockFtr.Length)
	if _, err = dst.Write(padding); err != nil {
		return err
	}
	if _, err = hasher.Write(padding); err != nil {
		return err
	}
	blockFtr.Checksum = hasher.Sum64()
	blockFtr.EOF = true
	return binary.Write(dst, binary.BigEndian, blockFtr)
}

func Receive(dst io.Writer, src io.Reader) error {

	var blockHdr _BlockHdr

	binary.Read(src, binary.BigEndian, &blockHdr)

	block := make([]byte, blockHdr.BlockSize)

	hasher := HashAlgorithm()

	r := io.TeeReader(src, hasher)

	var blockFtr _BlockFtr
	for {
		hasher.Reset()

		if _, err := io.ReadFull(r, block); err != nil {
			return ErrBrokenBlock
		}

		checksum := hasher.Sum64()

		if err := binary.Read(src, binary.BigEndian, &blockFtr); err != nil {
			return ErrBrokenBlock
		}

		if blockFtr.Checksum != checksum {
			return ErrBrokenBlock
		}

		if n, err := dst.Write(block[:blockFtr.Length]); err != nil || int64(n) != blockFtr.Length {
			return err
		}

		if blockFtr.EOF {
			break
		}
	}

	return nil
}
