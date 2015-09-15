// Command go-where finds where a Go identifier is defined and outputs
// the location in a format suitable for editor integration.
//
// Current values of GOOS, GOARCH etc are used.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tv42/where"
)

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s [-tags 'tag list'] IMPORTPATH[#IDENTIFIER[.FIELD]...[.METHOD]]\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Example:\n")
	fmt.Fprintf(os.Stderr, "  %s github.com/tv42/where#Ident\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func run(arg string, tags []string) error {
	buildCtx := build.Default
	buildCtx.BuildTags = tags
	pos, err := where.Ident(arg, &buildCtx)
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
	tags := flag.String("tags", "", "list of build tags, space-separated")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	arg := flag.Arg(0)

	taglist := strings.Fields(*tags)
	if err := run(arg, taglist); err != nil {
		log.Fatal(err)
	}
}
