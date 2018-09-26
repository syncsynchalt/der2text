package hinter

import (
	"github.com/syncsynchalt/der2text/read/indenter"
	"github.com/syncsynchalt/der2text/test"
	"strings"
	"testing"
)

func TestHinterPercents(t *testing.T) {
	test.Equals(t, false, isMostlyPrintable([]byte("")))
	test.Equals(t, true, isMostlyPrintable([]byte("abc")))
	test.Equals(t, true, isMostlyPrintable([]byte("abc\x00")))
	test.Equals(t, true, isMostlyPrintable([]byte("abc\x00\x01")))
	test.Equals(t, false, isMostlyPrintable([]byte("abc\x00\x01\x02")))
	test.Equals(t, true, isMostlyPrintable([]byte("abcdef\x00\x01\x02\x03\x04")))
}

func TestPrintHintNotPrintable(t *testing.T) {
	w := &strings.Builder{}
	ind := indenter.New(w)
	PrintHint(ind, []byte("\x10\x11\x12\x13"))
	test.Equals(t, "", w.String())
}

func TestPrintHintPrintable(t *testing.T) {
	w := &strings.Builder{}
	ind := indenter.New(w)
	PrintHint(ind, []byte("abc\x00def.\"g"))
	test.Equals(t, `# data: "abc.def..g"
`, w.String())
}

// helper function for the below tests
func hintTime(input string) string {
	w := &strings.Builder{}
	ind := indenter.New(w)
	PrintTimeHint(ind, []byte(input))
	return strings.TrimRight(w.String(), "\n")
}

func TestPrintTimeShort(t *testing.T) {
	test.Equals(t, "# 1969-12-31 00:11:22 GMT", hintTime("691231001122Z"))
	test.Equals(t, "# 1999-12-31 00:11:22 GMT", hintTime("991231001122Z"))
	test.Equals(t, "# 2018-09-10 01:02:03 GMT", hintTime("180910010203Z"))
	test.Equals(t, "# 2018-09-10 01:02:03", hintTime("180910010203"))
	test.Equals(t, "# 2018-09-10 01:02:03.123", hintTime("180910010203.123"))
	test.Equals(t, "# 2018-09-10 01:02:03.123 GMT", hintTime("180910010203.123Z"))
	test.Equals(t, "# 2018-09-10 01:02:03-0700", hintTime("180910010203-0700"))
	test.Equals(t, "# 2018-09-10 01:02:03+0700", hintTime("180910010203+0700"))
	test.Equals(t, "", hintTime("1x0910010203+0700"))
	test.Equals(t, "", hintTime("1809100102Z"))
}

func TestPrintTimeLong(t *testing.T) {
	test.Equals(t, "# 1999-12-31 00:11:22 GMT", hintTime("19991231001122Z"))
	test.Equals(t, "# 2018-09-10 01:02:03 GMT", hintTime("20180910010203Z"))
	test.Equals(t, "# 1918-09-10 01:02:03", hintTime("19180910010203"))
	test.Equals(t, "# 2018-09-10 01:02:03.123", hintTime("20180910010203.123"))
	test.Equals(t, "# 2018-09-10 01:02:03.123 GMT", hintTime("20180910010203.123Z"))
	test.Equals(t, "# 2018-09-10 01:02:03-0700", hintTime("20180910010203-0700"))
	test.Equals(t, "# 2018-09-10 01:02:03+0700", hintTime("20180910010203+0700"))
	test.Equals(t, "", hintTime("2018091001020304Z"))
}
