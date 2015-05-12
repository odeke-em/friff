package merkle

import (
	"fmt"
	"runtime"
	"testing"
)

func caller() string {
	_, thisFile, _, _ := runtime.Caller(1)
	return thisFile
}

func TestMd5Checksum(t *testing.T) {
	ckChan, err := checksumChanify(caller())
	if err != nil {
		t.Errorf("expected nil error got %v", err)
	}

	for ck := range ckChan {
		fmt.Println(ck)
	}
}
