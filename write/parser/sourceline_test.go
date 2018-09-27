package parser_test

import (
	"github.com/syncsynchalt/der2text/test"
	. "github.com/syncsynchalt/der2text/write/parser"
	"testing"
)

func TestSourceLineIndent(t *testing.T) {
	test.Equals(t, 0, SourceLine{Str: ""}.IndentLevel())
	test.Equals(t, 0, SourceLine{Str: "abc"}.IndentLevel())
	test.Equals(t, 0, SourceLine{Str: "abc  "}.IndentLevel())
	test.Equals(t, 1, SourceLine{Str: " abc"}.IndentLevel())
	test.Equals(t, 2, SourceLine{Str: "  abc"}.IndentLevel())
	test.Equals(t, 4, SourceLine{Str: "    abc"}.IndentLevel())

}

func TestSourceLineComment(t *testing.T) {
	test.Equals(t, true, SourceLine{Str: ""}.IsComment())
	test.Equals(t, false, SourceLine{Str: "a"}.IsComment())
	test.Equals(t, true, SourceLine{Str: "  "}.IsComment())
	test.Equals(t, false, SourceLine{Str: "  b"}.IsComment())
	test.Equals(t, true, SourceLine{Str: "  #"}.IsComment())
	test.Equals(t, true, SourceLine{Str: "  #  "}.IsComment())
	test.Equals(t, true, SourceLine{Str: "#  "}.IsComment())
	test.Equals(t, true, SourceLine{Str: "#"}.IsComment())
}

func TestSourceLineToken(t *testing.T) {
	expect := func(l *SourceLine, expected string) {
		test.CallerDepth = 2
		defer func() { test.CallerDepth = 1 }()

		word, err := l.NextToken()
		test.Ok(t, err)
		test.Equals(t, expected, word)
	}
	line := SourceLine{Str: ""}
	expect(&line, "")
	expect(&line, "")

	line = SourceLine{Str: "a b c"}
	expect(&line, "a")
	expect(&line, "b")
	expect(&line, "c")
	expect(&line, "")

	line = SourceLine{Str: " a b"}
	expect(&line, "a")
	expect(&line, "b")
	expect(&line, "")
}

func TestSourceLineTokenOctets(t *testing.T) {
	expect := func(l *SourceLine, expected string) {
		test.CallerDepth = 2
		defer func() { test.CallerDepth = 1 }()

		word, err := l.NextToken()
		test.Ok(t, err)
		test.Equals(t, expected, word)
	}
	line := SourceLine{Str: "foo :"}
	expect(&line, "foo")
	expect(&line, "")
	expect(&line, "")

	line = SourceLine{Str: "foo bar baz :010203"}
	expect(&line, "foo")
	expect(&line, "bar")
	expect(&line, "baz")
	expect(&line, "\x01\x02\x03")
	expect(&line, "")

	line = SourceLine{Str: ":010"}
	word, err := line.NextToken()
	test.Equals(t, "encoding/hex: odd length hex string", err.Error())
	test.Equals(t, "", word)
}

func TestSourceLineTokenString(t *testing.T) {
	expect := func(l *SourceLine, expected string) {
		test.CallerDepth = 2
		defer func() { test.CallerDepth = 1 }()

		word, err := l.NextToken()
		test.Ok(t, err)
		test.Equals(t, expected, word)
	}
	line := SourceLine{Str: "foo-bar '"}
	expect(&line, "foo-bar")
	expect(&line, "")
	expect(&line, "")

	line = SourceLine{Str: `foo-bar 'Hi\r\nmom`}
	expect(&line, "foo-bar")
	expect(&line, "Hi\r\nmom")
	expect(&line, "")
}

func TestSourceLineTokenType(t *testing.T) {
	line := SourceLine{Str: ""}
	test.Equals(t, "", line.NextTokenType())

	line = SourceLine{Str: "         "}
	test.Equals(t, "", line.NextTokenType())

	line = SourceLine{Str: "foo-bar"}
	test.Equals(t, "ATOM", line.NextTokenType())

	line = SourceLine{Str: "  foo-bar"}
	test.Equals(t, "ATOM", line.NextTokenType())

	line = SourceLine{Str: "'foo-bar"}
	test.Equals(t, "STRING", line.NextTokenType())

	line = SourceLine{Str: "      'foo-bar"}
	test.Equals(t, "STRING", line.NextTokenType())

	line = SourceLine{Str: ":0000"}
	test.Equals(t, "OCTETS", line.NextTokenType())

	line = SourceLine{Str: "      :"}
	test.Equals(t, "OCTETS", line.NextTokenType())
}

func TestSliceHigherIndentsEmpty(t *testing.T) {
	l := []SourceLine{}
	e := []SourceLine{}
	r := SliceHigherIndents(l, 10)
	test.Equals(t, e, r)
}

func TestSliceHigherIndentsEndsMid(t *testing.T) {
	l := []SourceLine{
		{Str: "  abc"},
		{Str: "  def"},
		{Str: "   gh"},
		{Str: "  abc"},
		{Str: " abc"},
	}
	e := []SourceLine{
		{Str: "  abc"},
		{Str: "  def"},
		{Str: "   gh"},
		{Str: "  abc"},
	}
	r := SliceHigherIndents(l, 1)
	test.Equals(t, e, r)
}

func TestSliceHigherIndentsNone(t *testing.T) {
	l := []SourceLine{
		{Str: "  abc"},
		{Str: "  def"},
	}
	e := []SourceLine{}
	r := SliceHigherIndents(l, 2)
	test.Equals(t, e, r)
}

func TestSliceHigherIndentsEndsEnd(t *testing.T) {
	l := []SourceLine{
		{Str: "  abc"},
		{Str: "  def"},
	}
	e := []SourceLine{
		{Str: "  abc"},
		{Str: "  def"},
	}
	r := SliceHigherIndents(l, 1)
	test.Equals(t, e, r)
}
