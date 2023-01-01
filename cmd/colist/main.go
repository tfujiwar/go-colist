package main

import (
	"fmt"
	"os"

	"github.com/tfujiwar/go-colist"
)

func main() {
	if err := colist.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
