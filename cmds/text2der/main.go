package main

import (
	"flag"
	"fmt"
	"github.com/syncsynchalt/der2text/write/parser"
	"os"
)

func main() {
	var err error
	calledHelp := flag.Bool("help", false, "This output")
	calledUsage := flag.Bool("usage", false, "This output")
	flag.Parse()
	if *calledHelp || *calledUsage || flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, `Usage: %s [input]

Parses [input] or stdin in the format output by der2text and
produces the PEM- or DER-encoded result on stdout`, os.Args[0])
		os.Exit(1)
	}

	in := os.Stdin
	if flag.NArg() == 1 {
		in, err = os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	}

	lines, err := parser.Lines(in)
	if err != nil {
		panic(err)
	}

	result, err := parser.Parse(lines)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(result)
}
