package parser_test

import (
	"github.com/syncsynchalt/der2text/test"
	. "github.com/syncsynchalt/der2text/write/parser"
	"strings"
	"testing"
)

func TestLinesCRLF(t *testing.T) {
	instr := "foo\r\nbar\r\n  baz\r\n  #bux\r\n"
	expected := []SourceLine{{Str: "foo"}, {Str: "bar"}, {Str: "  baz"}, {Str: "  #bux"}}
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	test.Equals(t, expected, lines)
}

func TestLinesLF(t *testing.T) {
	instr := "foo\nbar\n  baz\n  #bux\n"
	expected := []SourceLine{{Str: "foo"}, {Str: "bar"}, {Str: "  baz"}, {Str: "  #bux"}}
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	test.Equals(t, expected, lines)
}

func TestParsePem(t *testing.T) {
	instr := "PEM ENCODED X\n  UNIVERSAL PRIMITIVE INTEGER 0\n"
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)

	out, err := Parse(lines)
	test.Ok(t, err)
	test.Equals(t, "-----BEGIN X-----\nAgEA\n-----END X-----\n", string(out))
}
