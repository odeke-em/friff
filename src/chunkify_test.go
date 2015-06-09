package friff

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
	chunkMap, err := Chunkify(caller())
	if err != nil {
		t.Errorf("expected nil error got %v", err)
	}

	for id, shadow := range chunkMap {
		fmt.Println(id, shadow)
	}
}
