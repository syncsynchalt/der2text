package pem_test

import (
	"github.com/syncsynchalt/der2text/indenter"
	. "github.com/syncsynchalt/der2text/pem"
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

// helper function used by all tests below
func testPemData(tb testing.TB, input string, output string) {
	// run Parse, compare output
	var parseOut stringWriter
	out := indenter.New(&parseOut)
	err := Parse(out, []byte(input))
	if err != nil && err.Error() == output {
		return
	}
	test.Ok(tb, err)
	test.Equals(tb, output, parseOut.str)
}

func TestEmpty(t *testing.T) {
	testPemData(t, "", `Unable to parse PEM header`)
}

func TestPem(t *testing.T) {
	input := `-----BEGIN FOO BAR-----
AgEB
-----END FOO BAR-----
`
	expected := `PEM ENCODED FOO BAR
  UNIVERSAL PRIMITIVE INTEGER 1
`
	testPemData(t, input, expected)
}

func TestPemSpaces(t *testing.T) {
	// trailing spaces etc
	input := `
-----BEGIN FOO----- 

AgEB 

-----END FOO----- 
`
	expected := `PEM ENCODED FOO
  UNIVERSAL PRIMITIVE INTEGER 1
`
	testPemData(t, input, expected)
}

func TestPemCRLF(t *testing.T) {
	input := "-----BEGIN FOO-----\r\nAgEB\r\n-----END FOO-----\r\n"
	expected := `PEM ENCODED FOO
  UNIVERSAL PRIMITIVE INTEGER 1
`
	testPemData(t, input, expected)
}
