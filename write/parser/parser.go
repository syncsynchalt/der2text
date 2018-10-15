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
			lineno += 1
			continue
		}

		result, handled, err := parseLine(lines, lineno)
		if err != nil {
			return nil, fmt.Errorf("line %d: %s", lineno+1, err.Error())
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

	if len(typ) > 14 && typ[:14] == "UNHANDLED-TAG=" {
		if construction == "PRIMITIVE" {
			payload, err := line.NextToken()
			if err != nil {
				return nil, 0, err
			}
			data, err = der.WriteGeneric(class, construction, typ, payload)
			return data, 1, err
		} else if construction == "CONSTRUCTED" {
			innerLines := SliceHigherIndents(lines[lineno+1:], orig.IndentLevel())
			payload, err := Parse(innerLines)
			if err != nil {
				return nil, 0, err
			}
			data, err = der.WriteGeneric(class, construction, typ, string(payload))
			return data, len(innerLines) + 1, err
		} else {
			return nil, 0, fmt.Errorf("Unrecognized construction %s", construction)
		}
	}

	switch typ {
	case "END-OF-CONTENT":
		data, err := der.WriteEndOfContent(class, construction, typ)
		return data, 1, err
	case "BOOLEAN":
		value, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		data, err := der.WriteBoolean(class, construction, typ, value)
		return data, 1, err
	case "INTEGER", "ENUMERATED":
		tokenType := line.NextTokenType()
		token, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		if tokenType == "OCTETS" {
			data, err = der.WriteGeneric(class, construction, typ, token)
		} else {
			data, err = der.WriteInteger(class, construction, typ, token)
		}
		return data, 1, err
	case "BITSTRING":
		padding, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		payload, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		data, err = der.WriteBitstring(class, construction, typ, padding, payload)
		return data, 1, err
	case "NULL":
		data, err = der.WriteNull(class, construction, typ)
		return data, 1, err
	case "OID":
		oid, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		data, err = der.WriteOid(class, construction, typ, oid)
		return data, 1, err
	case "RELATIVEOID":
		oid, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		data, err = der.WriteRelativeOid(class, construction, typ, oid)
		return data, 1, err
	case "UNIVERSALSTRING":
		str, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		data, err = der.WriteUniversalString(class, construction, typ, str)
		return data, 1, err
	case "BMPSTRING":
		str, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		data, err = der.WriteBMPString(class, construction, typ, str)
		return data, 1, err
	case "OCTETSTRING", "OBJECTDESCRIPTION", "EXTERNAL", "REAL", "EMBEDDED-PDV", "UTF8STRING",
		"NUMERICSTRING", "PRINTABLESTRING", "T61STRING", "VIDEOTEXSTRING", "IA5STRING", "UTCTIME",
		"GENERALIZEDTIME", "GRAPHICSTRING", "VISIBLESTRING", "GENERALSTRING", "CHARACTERSTRING":
		// no special handling on these types, they're just bare payload for now
		payload, err := line.NextToken()
		if err != nil {
			return nil, 0, err
		}
		data, err = der.WriteGeneric(class, construction, typ, payload)
		return data, 1, err
	case "SET", "SEQUENCE":
		innerLines := SliceHigherIndents(lines[lineno+1:], orig.IndentLevel())
		payload, err := Parse(innerLines)
		if err != nil {
			return nil, 0, err
		}
		data, err = der.WriteGeneric(class, construction, typ, string(payload))
		return data, len(innerLines) + 1, err
	default:
		return nil, 0, fmt.Errorf("Unrecognized type %s", typ)
	}
}
