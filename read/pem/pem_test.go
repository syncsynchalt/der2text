package pem_test

import (
	"github.com/syncsynchalt/der2text/read/indenter"
	. "github.com/syncsynchalt/der2text/read/pem"
	"github.com/syncsynchalt/der2text/test"
	"strings"
	"testing"
)

// helper function used by all tests below
func testPemData(tb testing.TB, input string, output string) {
	test.CallerDepth = 2
	defer func() { test.CallerDepth = 1 }()

	// run Parse, compare output
	var parseOut strings.Builder
	out := indenter.New(&parseOut)
	err := Parse(out, []byte(input))
	if err != nil && err.Error() == output {
		return
	}
	test.Ok(tb, err)
	test.Equals(tb, output, parseOut.String())
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
