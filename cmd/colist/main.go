package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tfujiwar/go-colist"
)

func main() {
	log.SetFlags(0)

	if err := colist.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
