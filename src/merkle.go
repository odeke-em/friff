package merkle

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type chunk struct {
	id            int64
	data          []byte
	md5Checksum   string
	cacheChecksum bool
}

func (ck *chunk) Read(p []byte) (int, error) {
	return copy(ck.data, p), io.EOF
}

func chunks(absPath string, chunkSize int64) ([]chunk, error) {
	lInfo, lErr := os.Stat(absPath)
	if lErr != nil || lInfo == nil {
		return nil, lErr
	}

	body, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	chunkCount := 1 + (lInfo.Size() / chunkSize)
	chunkList := make([]chunk, chunkCount)

	i, n := int64(0), int(0)
	for {
		buf := make([]byte, chunkSize)
		n, err = io.ReadFull(body, buf)
		buf = buf[:n]
		chunkList[i] = chunk{
			cacheChecksum: true,
			data:          buf[:n],
			id:            i,
		}
		i += 1

		if err != nil {
			break
		}
	}

	chunkList = chunkList[:i]
	return chunkList, nil
}

func Chunks(absPath string) ([]chunk, error) {
	return chunks(absPath, 4*1024*1024)
}

func (ck *chunk) String() string {
	return md5Checksum(ck)
}

func md5Checksum(ck *chunk) string {
	if ck.md5Checksum != "" && ck.cacheChecksum {
		return ck.md5Checksum
	}
	ck.md5Checksum = fmt.Sprintf("%x", md5.Sum(ck.data))
	return ck.md5Checksum
}
