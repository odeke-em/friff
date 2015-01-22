package main

import (
	"fmt"
	"os"

	merkle "github.com/odeke-em/merkle-go/src"
)

func main() {
	chunks, _ := merkle.Chunks(os.Args[1])
	for i, chunk := range chunks {
		fmt.Println(i, chunk.String())
	}
}
