package main

import (
	"flag"
	"fmt"
	"github.com/syncsynchalt/der2text/der"
	"github.com/syncsynchalt/der2text/indenter"
	"github.com/syncsynchalt/der2text/pem"
	"io/ioutil"
	"os"
)

func main() {
	var err error
	calledHelp := flag.Bool("help", false, "This output")
	calledUsage := flag.Bool("usage", false, "This output")
	flag.Parse()
	if *calledHelp || *calledUsage || flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, `Usage: %s [input]

Parses [input] or stdin as a DER-encoded or PEM-encoded file and
produces a more readable output on stdout`, os.Args[0])
		os.Exit(1)
	}

	in := os.Stdin
	if flag.NArg() == 1 {
		in, err = os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	}

	data, err := ioutil.ReadAll(in)
	if err != nil {
		panic(err)
	}

	if len(data) > 11 && string(data[:11]) == "-----BEGIN " {
		pem.Parse(indenter.New(os.Stdout), data)
	} else {
		der.Parse(indenter.New(os.Stdout), data)
	}
}
