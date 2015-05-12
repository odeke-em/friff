package merkle

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"syscall"
)

var (
	KB            = 1024
	BytesPerBlock = 4 * KB
)

type chunk struct {
	id   uint
	data []byte
	n    int
}

type Shadow struct {
	id       uint
	size     int
	checksum string
}

func md5Checksum(bst []byte) string {
	return fmt.Sprintf("%x", md5.Sum(bst))
}

func (ck *chunk) compute() chan *Shadow {
	done := make(chan *Shadow)

	go func() {
		shad := Shadow{
			id:       ck.id,
			size:     ck.n,
			checksum: md5Checksum(ck.data),
		}

		done <- &shad
		close(done)
	}()

	return done
}

func chunkFile(blobAt string) (ckl chan *chunk, err error) {
	var fh *os.File
	fh, err = os.Open(blobAt)
	if err != nil {
		return
	}

	ckl = make(chan *chunk)
	go func() {
		defer fh.Close()
		i := uint(0)
		for {
			bts := make([]byte, BytesPerBlock)
			n, err := io.ReadAtLeast(fh, bts, 1)

			if err != nil {
				if err == io.EOF {
					break
				} else if err == syscall.EINTR {
					continue
				} else {
					fmt.Printf("error on ith block: %d err: %v\n", i, err)
					break
				}
			}

			ckl <- &chunk{id: i, data: bts, n: n}
			i += 1
		}
		close(ckl)
	}()

	return
}

func checksumChanify(blobAt string) (chan *Shadow, error) {
	ckl, err := chunkFile(blobAt)
	if err != nil {
		return nil, err
	}

	ckll := make(chan *Shadow)
	go func() {
		chanOChan := make(chan chan *Shadow)

		go func() {
			for cch := range ckl {
				chanOChan <- cch.compute()
			}
			close(chanOChan)
		}()

		for cksumChan := range chanOChan {
			cksum := <-cksumChan
			ckll <- cksum
		}
		close(ckll)
	}()

	return ckll, nil
}
