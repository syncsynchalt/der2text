package parser

import (
	"bufio"
	"fmt"
	"github.com/syncsynchalt/der2text/write/der"
	"github.com/syncsynchalt/der2text/write/pem"
	"io"
	"strings"
)

func Lines(in io.Reader) ([]SourceLine, error) {
	var inlines []SourceLine
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, "\r\n")
		inlines = append(inlines, SourceLine{Str: line})
	}
	err := scanner.Err()
	if err != nil {
		return []SourceLine{}, err
	}

	return inlines, nil
}

// decorate error in line number
func Parse(lines []SourceLine) ([]byte, error) {
	out := make([]byte, 0)
	lineno := 0
	for lineno < len(lines) {
		line := lines[lineno]
		if line.IsComment() {
			continue
		}

		result, handled, err := parseLine(lines, lineno)
		if err != nil {
			return nil, fmt.Errorf("line %d: %s", lineno, err.Error())
		}
		out = append(out, result...)
		lineno += handled
	}
	return out, nil
}

func parseLine(lines []SourceLine, lineno int) (result []byte, handled int, err error) {
	line := lines[lineno]
	orig := line

	class, err := line.NextToken()
	if err != nil {
		return nil, 0, err
	}
	if class == "PEM" {
		word2, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		if word2 != "ENCODED" {
			return nil, 0, fmt.Errorf("Incorrent PEM line format %s", orig.Str)
		}
		innerlines := lines[lineno+1:]
		data, err := Parse(innerlines)
		if err != nil {
			return nil, 0, err
		}
		return pem.PemEncode(line.Str, data), len(lines), nil
	}
	construction, err := line.NextToken()
	if err != nil {
		return nil, 0, err
	}
	typ, err := line.NextToken()
	if err != nil {
		return nil, 0, err
	}

	var data []byte
	switch typ {
	case "INTEGER", "ENUMERATED":
		tokenType := line.NextTokenType()
		token, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		if tokenType == "OCTETS" {
			data, err = der.WriteIntegerPreserved(class, construction, typ, token)
		} else {
			data, err = der.WriteInteger(class, construction, typ, token)
		}
		return data, 1, err
	default:
		return nil, 0, fmt.Errorf("Unrecognized type %s", typ)
	}
}
