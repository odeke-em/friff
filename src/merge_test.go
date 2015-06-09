package friff

import (
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {
	chunkMap, err := Chunkify(caller())
	if err != nil {
		t.Errorf("expected nil error got %v", err)
	}

	for id, shadow := range chunkMap {
		fmt.Println(id, shadow)
	}
}
