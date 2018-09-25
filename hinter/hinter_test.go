package hinter

import (
	"github.com/syncsynchalt/der2text/indenter"
	"github.com/syncsynchalt/der2text/test"
	"testing"
)

// an io.Writer that builds a string
type stringWriter struct {
	str string
}

func (s *stringWriter) Write(p []byte) (n int, err error) {
	s.str += string(p)
	return len(p), nil
}

func TestHinterPercents(t *testing.T) {
	test.Equals(t, false, isMostlyPrintable([]byte("")))
	test.Equals(t, true, isMostlyPrintable([]byte("abc")))
	test.Equals(t, true, isMostlyPrintable([]byte("abc\x00")))
	test.Equals(t, true, isMostlyPrintable([]byte("abc\x00\x01")))
	test.Equals(t, false, isMostlyPrintable([]byte("abc\x00\x01\x02")))
	test.Equals(t, true, isMostlyPrintable([]byte("abcdef\x00\x01\x02\x03\x04")))
}

func TestPrintHintNotPrintable(t *testing.T) {
	w := &stringWriter{}
	ind := indenter.New(w)
	PrintHint(ind, []byte("\x10\x11\x12\x13"))
	test.Equals(t, "", w.str)
}

func TestPrintHintPrintable(t *testing.T) {
	w := &stringWriter{}
	ind := indenter.New(w)
	PrintHint(ind, []byte("abc\x00def.\"g"))
	test.Equals(t, `# data: "abc.def..g"
`, w.str)
}
