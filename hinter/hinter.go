package hinter

import (
	"github.com/syncsynchalt/der2text/printer"
)

func isMostlyPrintable(content []byte) bool {
	if len(content) < 1 {
		return false
	}
	printable := 0.0
	for _, v := range content {
		if v >= 0x20 && v < 0x7f {
			printable++
		}
	}
	if len(content) > 0 && (printable/float64(len(content))) > 0.5 {
		return true
	}
	return false
}

func PrintHint(out printer.Printer, content []byte) {
	if isMostlyPrintable(content) {
		out.Print("# data: \"")
		for _, v := range content {
			if v == '"' || v < 0x20 || v >= 0x7f {
				out.Print(".")
			} else {
				out.Printf("%c", v)
			}
		}
		out.Print("\"\n")
	}
}
