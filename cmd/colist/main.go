package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tfujiwar/go-colist"
)

func main() {
	var output string
	flag.StringVar(&output, "o", "text", "output format")
	flag.StringVar(&output, "output", "text", "output format")

	var dir string
	flag.StringVar(&dir, "d", ".", "repository directory")
	flag.StringVar(&dir, "dir", ".", "repository directory")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "show debug log")
	flag.BoolVar(&verbose, "verbose", false, "show debug log")

	var help bool
	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&help, "help", false, "show help")

	flag.Parse()

	if help {
		usage(os.Stdout)
		os.Exit(0)
	}

	if verbose {
		log.SetFlags(0)
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}

	var formatFunc func([]*colist.Rule, io.Writer) error
	switch output {
	case "text":
		formatFunc = colist.OutputText
	case "json":
		formatFunc = colist.OutputJson
	default:
		fmt.Fprintf(os.Stderr, "[ERROR] not supported output: select \"text\" or \"json\"\n")
		os.Exit(1)
	}

	args := flag.Args()

	var remote string
	var baseBranch string
	switch len(args) {
	case 0:
		remote = ""
		baseBranch = ""
	case 1:
		remote = ""
		baseBranch = args[0]
	case 2:
		remote = args[0]
		baseBranch = args[1]
	default:
		usage(os.Stderr)
		os.Exit(1)
	}

	log.Printf("[DEBUG] output     : %v\n", output)
	log.Printf("[DEBUG] dir        : %v\n", dir)
	log.Printf("[DEBUG] verbose    : %v\n", verbose)
	log.Printf("[DEBUG] help       : %v\n", help)
	log.Printf("[DEBUG] remote     : %v\n", remote)
	log.Printf("[DEBUG] baseBranch : %v\n", baseBranch)

	rules, err := colist.Main(dir, remote, baseBranch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	err = formatFunc(rules, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "List GitHub CODEOWNERS of changed files on a current branch\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Usage:\n")
	fmt.Fprintf(w, "  colist [flags]                        : compare with remote or local main branch\n")
	fmt.Fprintf(w, "  colist [flags] <BASE_BRANCH>          : compare with remote or local <BASE_BRANCH>\n")
	fmt.Fprintf(w, "  colist [flags] <REMOTE> <BASE_BRANCH> : compare with <REMOTE>/<BASE_BRANCH>\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Flags:\n")
	fmt.Fprintf(w, "  -o, --output text|json : output format\n")
	fmt.Fprintf(w, "  -d, --dir <DIR>        : repository directory\n")
	fmt.Fprintf(w, "  -v, --verbose          : show debug log\n")
	fmt.Fprintf(w, "  -h, --help             : show this message\n")
	fmt.Fprintf(w, "\n")
}
