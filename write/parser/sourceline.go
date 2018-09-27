package parser

import (
	"encoding/hex"
	"strings"
)

type SourceLine struct {
	Str string
}

func (s SourceLine) IndentLevel() int {
	return len(s.Str) - len(strings.TrimLeft(s.Str, " "))
}

func (s SourceLine) IsComment() bool {
	str := strings.Trim(s.Str, " \t\r\n")
	return len(str) == 0 || str[0] == '#'
}

// note: this modifies the line and invalidates IndentLevel() and IsComment()
func (s *SourceLine) NextToken() (string, error) {
	// skip past any spaces
	s.Str = strings.TrimLeft(s.Str, " ")

	if s.Str == "" {
		return "", nil
	}

	// either we have :octets, 'string, or an atom delim'd by space
	if s.Str[0] == ':' {
		// rest of line is octets in hex
		data, err := hex.DecodeString(s.Str[1:])
		if err != nil {
			return "", err
		}
		s.Str = ""
		return string(data), nil
	} else if s.Str[0] == '\'' {
		// rest of line is string with \r\n escaped
		rest := s.Str[1:]
		s.Str = ""
		rest = strings.Replace(rest, "\\r", "\r", -1)
		rest = strings.Replace(rest, "\\n", "\n", -1)
		return rest, nil
	} else {
		// grab the space-delimited word
		i := strings.IndexRune(s.Str, ' ')
		if i < 0 {
			i = len(s.Str)
		}
		word := s.Str[:i]
		s.Str = s.Str[i:]
		s.Str = strings.TrimLeft(s.Str, " ")
		return word, nil
	}
}

func (s SourceLine) NextTokenType() string {
	ss := strings.TrimLeft(s.Str, " ")
	if len(ss) == 0 {
		return ""
	}
	switch ss[0] {
	case ':':
		return "OCTETS"
	case '\'':
		return "STRING"
	default:
		return "ATOM"
	}
}

// return slice of the elements at the head of lines that have an indent level higher than indent
func SliceHigherIndents(lines []SourceLine, indent int) []SourceLine {
	i := 0
	for i < len(lines) && lines[i].IndentLevel() > indent {
		i++
	}
	return lines[:i]
}
