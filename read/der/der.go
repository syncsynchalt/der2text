package der

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/syncsynchalt/der2text/read/hinter"
	"github.com/syncsynchalt/der2text/read/indenter"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"strconv"
)

const (
	classUniversal       = 0 << 6
	classApplication     = 1 << 6
	classContextSpecific = 2 << 6
	classPrivate         = 3 << 6

	constructed = 1 << 5
	primitive   = 0 << 5

	typeEndOfContent      = 0x0
	typeBoolean           = 0x1
	typeInteger           = 0x2
	typeBitString         = 0x3
	typeOctetString       = 0x4
	typeNull              = 0x5
	typeObjectIdentifier  = 0x6
	typeObjectDescription = 0x7
	typeExternal          = 0x8
	typeReal              = 0x9
	typeEnumerated        = 0xA
	typeEmbeddedPDV       = 0xB
	typeUtf8String        = 0xC
	typeRelativeOID       = 0xD
	typeSequence          = 0x10
	typeSet               = 0x11
	typeNumericString     = 0x12
	typePrintableString   = 0x13
	typeT61String         = 0x14
	typeVideotexString    = 0x15
	typeIA5String         = 0x16
	typeUTCTime           = 0x17
	typeGeneralizedTime   = 0x18
	typeGraphicString     = 0x19
	typeVisibleString     = 0x1A
	typeGeneralString     = 0x1B
	typeUniversalString   = 0x1C
	typeCharacterString   = 0x1D
	typeBMPString         = 0x1E
	typeIsLongFormTag     = 0x1F
)

func Parse(out *indenter.Indenter, data []byte) error {
	for len(data) > 0 {
		var err error
		data, err = parseElement(out, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseElement(out *indenter.Indenter, data []byte) (rest []byte, err error) {
	if len(data) < 2 {
		return nil, errors.New("short DER read, need at least two bytes, got " + strconv.Itoa(len(data)))
	}

	typeByte := data[0]
	typeTag := int(typeByte & 0x1F)
	typeConstructed := typeByte & 0x20
	typeClass := typeByte & 0xC0
	numBytesInType := 1

	if typeTag == typeIsLongFormTag {
		typeTag, numBytesInType, err = readLongType(data)
		if err != nil {
			return nil, err
		}
	}

	switch typeClass {
	case classUniversal:
		out.Print("UNIVERSAL ")
	case classApplication:
		out.Print("APPLICATION ")
	case classContextSpecific:
		out.Print("CONTEXT-SPECIFIC ")
	case classPrivate:
		out.Print("PRIVATE ")
	}

	switch typeConstructed {
	case primitive:
		out.Print("PRIMITIVE ")
	case constructed:
		out.Print("CONSTRUCTED ")
	}

	contentLen, rest, err := decodeLength(data[numBytesInType:])
	if err != nil {
		return nil, err
	}
	if len(rest) < contentLen {
		return nil, errors.New("Short content, need " + strconv.Itoa(contentLen) +
			" bytes but have " + strconv.Itoa(len(rest)))
	}
	content := rest[:contentLen]
	rest = rest[contentLen:]

	switch typeByte {
	case typeEndOfContent | primitive:
		if contentLen != 0 {
			return nil, errors.New("End-of-content had unexpected length " + strconv.Itoa(contentLen))
		}
		out.Println("END-OF-CONTENT")
	case typeBoolean | primitive:
		if contentLen != 1 {
			return nil, errors.New("Boolean had unexpected length " + strconv.Itoa(contentLen))
		}
		if content[0] == byte(0) {
			out.Println("BOOLEAN FALSE")
		} else {
			out.Println("BOOLEAN TRUE")
		}
	case typeInteger | primitive:
		handleInteger("INTEGER", out, content)
	case typeBitString | primitive:
		if contentLen < 1 {
			return nil, errors.New("BitString had no padding byte")
		}
		padding := int(content[0])
		if padding < 0 || padding > 7 {
			return nil, errors.New("BitString padding has illegal value " + strconv.Itoa(padding))
		}
		out.Printf("BITSTRING PAD=%d ", padding)
		printOctets(out, content[1:])
		out.Print("\n")
	case typeOctetString | primitive:
		handleData("OCTETSTRING", out, content)
	case typeNull | primitive:
		if contentLen != 0 {
			return nil, errors.New("Null has non-zero content")
		}
		out.Print("NULL\n")
	case typeObjectIdentifier | primitive:
		if contentLen < 1 {
			return nil, errors.New("OID doesn't have content")
		}
		first := content[0] / 40
		second := content[0] % 40
		oid := fmt.Sprintf("%d.%d", first, second)
		var build int
		for _, v := range content[1:] {
			build *= 128
			build += int(v & 0x7f)
			if v&0x80 == 0 {
				oid += fmt.Sprintf(".%d", build)
				build = 0
			}
		}
		out.Println("OID", oid)
		hinter.PrintOidName(out, oid)
	case typeObjectDescription | primitive:
		handleData("OBJECTDESCRIPTION", out, content)
	case typeExternal | constructed:
		handleData("EXTERNAL", out, content)
	case typeReal | primitive:
		handleData("REAL", out, content)
	case typeEnumerated | primitive:
		handleInteger("ENUMERATED", out, content)
	case typeEmbeddedPDV | constructed:
		handleData("EMBEDDED-PDV", out, content)
	case typeUtf8String | primitive:
		handleString("UTF8STRING", out, content)
	case typeRelativeOID | primitive:
		if contentLen < 1 {
			return nil, errors.New("Relative OID doesn't have content")
		}
		oid := ""
		var build int
		for _, v := range content {
			build *= 128
			build += int(v & 0x7f)
			if v&0x80 == 0 {
				oid += fmt.Sprintf(".%d", build)
				build = 0
			}
		}
		oid = oid[1:]
		out.Println("RELATIVEOID", oid)
		hinter.PrintOidName(out, oid)
	case typeNumericString | primitive:
		handleString("NUMERICSTRING", out, content)
	case typePrintableString | primitive:
		handleString("PRINTABLESTRING", out, content)
	case typeSet | constructed:
		out.Println("SET")
		Parse(out.NextLevel(), content)
	case typeSequence | constructed:
		out.Println("SEQUENCE")
		Parse(out.NextLevel(), content)
	case typeT61String | primitive:
		// handleString might be fine? just needs to be round-trip safe
		handleData("T61STRING", out, content)
	case typeVideotexString | primitive:
		// handleString might be fine? just needs to be round-trip safe
		handleData("VIDEOTEXSTRING", out, content)
	case typeIA5String | primitive:
		handleString("IA5STRING", out, content)
	case typeUTCTime | primitive:
		handleString("UTCTIME", out, content)
		hinter.PrintTimeHint(out, content)
	case typeGeneralizedTime | primitive:
		handleString("GENERALIZEDTIME", out, content)
		hinter.PrintTimeHint(out, content)
	case typeGraphicString | primitive:
		// handleString might be fine? just needs to be round-trip safe
		handleData("GRAPHICSTRING", out, content)
	case typeVisibleString | primitive:
		handleString("VISIBLESTRING", out, content)
	case typeGeneralString | primitive:
		// handleString might be fine? just needs to be round-trip safe
		handleData("GENERALSTRING", out, content)
	case typeUniversalString | primitive:
		b, err := utf32ToUtf8(content)
		if err != nil {
			return nil, err
		}
		handleString("UNIVERSALSTRING", out, b)
	case typeCharacterString | primitive:
		// handleString might be fine? just needs to be round-trip safe
		handleData("CHARACTERSTRING", out, content)
	case typeBMPString | primitive:
		b, err := utf16ToUtf8(content)
		if err != nil {
			return nil, err
		}
		handleString("BMPSTRING", out, b)
	default:
		if typeConstructed == constructed {
			out.Printf("UNHANDLED-TAG=%d\n", typeTag)
			Parse(out.NextLevel(), content)
		} else {
			label := fmt.Sprintf("UNHANDLED-TAG=%d", typeTag)
			handleData(label, out, content)
		}
	}

	return rest, nil
}

func readLongType(data []byte) (tag int, numBytesEncoded int, err error) {
	if data[0]&0x1F != 0x1F {
		panic("readLongType was passed a non-long type")
	}
	var result, i int
	for i = 1; i < len(data); i++ {
		result = result << 7
		result += int(data[i] & 0x7f)
		if data[i]&0x80 == 0 {
			return result, i + 1, nil
		}
	}
	return 0, 0, errors.New("ran out of bytes while reading long type tag")
}

func decodeLength(data []byte) (length int, rest []byte, err error) {
	firstByte := data[0]
	if firstByte&0x80 != 0 {
		numToRead := int(firstByte ^ 0x80)
		if len(data)-1 < numToRead {
			return 0, []byte{}, errors.New("Can't satisfy request to read " +
				strconv.Itoa(numToRead) + " bytes to get length")
		}
		length := 0
		for i := 0; i < numToRead; i++ {
			length *= 256
			length += int(data[1+i])
		}
		return length, data[1+numToRead:], nil
	} else {
		return int(firstByte), data[1:], nil
	}
}

func printString(out *indenter.Indenter, content []byte) {
	bout := bufio.NewWriter(out)
	for _, v := range content {
		if v == '\n' {
			bout.Write([]byte("\\n"))
		} else if v == '\r' {
			bout.Write([]byte("\\r"))
		} else {
			bout.Write([]byte{v})
		}
	}
	bout.Flush()
}

func handleData(label string, out *indenter.Indenter, content []byte) {
	out.Printf("%s ", label)
	printOctets(out, content)
	out.Print("\n")
	hinter.PrintStrings(out, content)
}

func handleString(label string, out *indenter.Indenter, content []byte) {
	out.Printf("%s '", label)
	printString(out, content)
	out.Print("\n")
}

func handleInteger(label string, out *indenter.Indenter, content []byte) {
	if len(content) > 0 && len(content) < 8 && content[0]&0x80 == 0 {
		// conveniently display it
		value := int64(0)
		if content[0]&0x80 == 0 {
			// positive number
			for _, v := range content {
				value *= 256
				value += int64(v)
			}
		}
		out.Println(label, value)
	} else if len(content) >= 8 || len(content) == 0 || content[0]&0x80 != 0 {
		// just dump it in hex
		handleData(label, out, content)
	}
}

func printOctets(out *indenter.Indenter, content []byte) {
	out.Print(":")
	for _, v := range content {
		out.Printf("%02X", v)
	}
}

func utf16ToUtf8(input []byte) ([]byte, error) {
	decoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
	return decoder.Bytes(input)
}

func utf32ToUtf8(input []byte) ([]byte, error) {
	decoder := utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewDecoder()
	return decoder.Bytes(input)
}
