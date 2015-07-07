package friff

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"syscall"
)

var (
	KB            = 1024
	BytesPerChunk = 256 * KB
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

func noop() {}

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

	fhClose := func() {
		fh.Close()
	}

	return chunkChaner(fh, fhClose)
}

func chunkChaner(fh io.Reader, deferal func()) (ckl chan *chunk, err error) {
	ckl = make(chan *chunk)

	go func() {
		defer func() {
			close(ckl)
			deferal()
		}()

		i := uint(0)
		for {
			bts := make([]byte, BytesPerChunk)
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
			defer close(chanOChan)
			for cch := range ckl {
				chanOChan <- cch.compute()
			}
		}()

		for cksumChan := range chanOChan {
			cksum := <-cksumChan
			ckll <- cksum

			// fmt.Println("cksum", cksum)
		}
		close(ckll)
	}()

	return ckll, nil
}

func Chunkify(blobAt string) (map[uint]*Shadow, error) {
	chunkChan, err := checksumChanify(blobAt)
	if err != nil {
		return nil, err
	}

	results := make(map[uint]*Shadow)
	for shad := range chunkChan {
		results[shad.id] = shad
	}

	return results, nil
}
