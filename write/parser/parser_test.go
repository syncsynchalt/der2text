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

// just breaking these end-to-end tests in workable chunks
func TestParseTypes1(t *testing.T) {
	instr := `#
UNIVERSAL PRIMITIVE END-OF-CONTENT
UNIVERSAL PRIMITIVE BOOLEAN TRUE
UNIVERSAL PRIMITIVE INTEGER -1
UNIVERSAL PRIMITIVE BITSTRING PAD=1 :A0
UNIVERSAL PRIMITIVE OCTETSTRING :010203
UNIVERSAL PRIMITIVE NULL
`
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	out, err := Parse(lines)
	test.Ok(t, err)
	test.Equals(t, []byte("\x00\x00\x01\x01\xff\x02\x01\xFF\x03\x02\x01\xA0\x04\x03\x01\x02\x03\x05\x00"), out)
}

func TestParseTypes2(t *testing.T) {
	instr := `#
UNIVERSAL PRIMITIVE OID 1.2.3
UNIVERSAL PRIMITIVE OBJECTDESCRIPTION :0102
UNIVERSAL CONSTRUCTED EXTERNAL :030405
UNIVERSAL PRIMITIVE REAL :03312E452B30
`
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	out, err := Parse(lines)
	test.Ok(t, err)
	test.Equals(t, []byte("\x06\x02\x2A\x03\x07\x02\x01\x02\x28\x03\x03\x04\x05\x09\x06\x031.E+0"), out)
}

func TestParseTypes3(t *testing.T) {
	instr := `#
UNIVERSAL PRIMITIVE ENUMERATED :01
UNIVERSAL CONSTRUCTED EMBEDDED-PDV :4142
UNIVERSAL PRIMITIVE UTF8STRING 'hí mom
UNIVERSAL PRIMITIVE RELATIVEOID 2.3
`
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	out, err := Parse(lines)
	test.Ok(t, err)
	test.Equals(t, []byte("\x0A\x01\x01\x2B\x02\x41\x42\x0C\x07h\xC3\xAD mom\x0D\x02\x02\x03"), out)
}

func TestParseTypes4(t *testing.T) {
	instr := `#
UNIVERSAL PRIMITIVE NUMERICSTRING 'a
UNIVERSAL PRIMITIVE PRINTABLESTRING 'b
UNIVERSAL PRIMITIVE T61STRING 'c
UNIVERSAL PRIMITIVE VIDEOTEXSTRING 'd
UNIVERSAL PRIMITIVE IA5STRING 'e
UNIVERSAL PRIMITIVE UTCTIME 'f
`
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	out, err := Parse(lines)
	test.Ok(t, err)
	test.Equals(t, []byte("\x12\x01\x61\x13\x01\x62\x14\x01\x63\x15\x01\x64\x16\x01\x65\x17\x01\x66"), out)
}

func TestParseTypes5(t *testing.T) {
	instr := `#
UNIVERSAL PRIMITIVE GENERALIZEDTIME 'a
UNIVERSAL PRIMITIVE GRAPHICSTRING 'b
UNIVERSAL PRIMITIVE VISIBLESTRING 'c
UNIVERSAL PRIMITIVE GENERALSTRING 'd
UNIVERSAL PRIMITIVE CHARACTERSTRING 'e
PRIVATE PRIMITIVE UNHANDLED-TAG=20 'f
`
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	out, err := Parse(lines)
	test.Ok(t, err)
	test.Equals(t, []byte("\x18\x01\x61\x19\x01\x62\x1A\x01\x63\x1B\x01\x64\x1D\x01\x65\xD4\x01\x66"), out)
}

func TestParseTypes6(t *testing.T) {
	instr := `#
UNIVERSAL PRIMITIVE UNIVERSALSTRING 'a
UNIVERSAL PRIMITIVE BMPSTRING 'b
`
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	out, err := Parse(lines)
	test.Ok(t, err)
	test.Equals(t, []byte("\x1C\x04\x00\x00\x00a\x1E\x02\x00b"), out)
}

func TestParseTypesSetSeq(t *testing.T) {
	instr := `#
UNIVERSAL CONSTRUCTED SET
  UNIVERSAL CONSTRUCTED SEQUENCE
    UNIVERSAL PRIMITIVE NULL
  UNIVERSAL PRIMITIVE END-OF-CONTENT
UNIVERSAL PRIMITIVE INTEGER 1
`
	lines, err := Lines(strings.NewReader(instr))
	test.Ok(t, err)
	out, err := Parse(lines)
	test.Ok(t, err)
	test.Equals(t, []byte("\x31\x06\x30\x02\x05\x00\x00\x00\x02\x01\x01"), out)
}
