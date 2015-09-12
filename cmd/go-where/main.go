// Command go-where finds where a Go identifier is defined and outputs
// the location in a format suitable for editor integration.
//
// Current values of GOOS, GOARCH etc are used.
//
// TODO support -tags
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tv42/where"
)

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s IMPORTPATH[#IDENTIFIER[.FIELD]...[.METHOD]]\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Example:\n")
	fmt.Fprintf(os.Stderr, "  %s github.com/tv42/where#Ident\n", prog)
}

func run(arg string) error {
	pos, err := where.Ident(arg, nil)
	if err != nil {
		return err
	}
	if _, err := fmt.Println(pos); err != nil {
		return err
	}
	return nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	arg := flag.Arg(0)

	if err := run(arg); err != nil {
		log.Fatal(err)
	}
}
