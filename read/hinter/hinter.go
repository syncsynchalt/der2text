package hinter

// generally useful human-readable hints for what the data might represent

import (
	"github.com/syncsynchalt/der2text/read/indenter"
	"strconv"
	"strings"
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

// if the data looks like it's mostly text then print what we can as a hint to the human reading it
func PrintHint(out *indenter.Indenter, content []byte) {
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

func PrintTimeHint(out *indenter.Indenter, content []byte) {
	// there are many mutually ambiguous representations,
	// see https://www.obj-sys.com/asn1tutorial/node14.html if you want
	// more examples to implement

	s := string(content)
	tz := ""
	if len(s) > 0 && s[len(s)-1] == 'Z' {
		s = s[:len(s)-1]
		tz = " GMT"
	} else if i := strings.IndexAny(s, "+-"); i != -1 {
		tz = s[i:]
		s = s[:i]
	}

	msecs := ""
	if i := strings.IndexRune(s, '.'); i != -1 {
		msecs = s[i:]
		s = s[:i]
	}

	var y, m, d, hh, mm, ss string
	if len(s) == 14 {
		y, m, d, hh, mm, ss = s[0:4], s[4:6], s[6:8], s[8:10], s[10:12], s[12:14]
	} else if len(s) == 12 {
		y, m, d, hh, mm, ss = s[0:2], s[2:4], s[4:6], s[6:8], s[8:10], s[10:12]
		tmp, err := strconv.Atoi(y)
		if err != nil {
			return
		}
		if tmp < 69 {
			y = "20" + y
		} else {
			y = "19" + y
		}
	} else {
		// ambiguous, not going to guess
		return
	}
	out.Printf("# %s-%s-%s %s:%s:%s%s%s\n", y, m, d, hh, mm, ss, msecs, tz)
}
