package der

import (
	"github.com/syncsynchalt/der2text/indenter"
	"encoding/base64"
	"strconv"
)

const (
	classUniversal       = 0 << 6
	classApplication     = 1 << 6
	classContextSpecific = 2 << 6
	classPrivate         = 3 << 6

	composed  = 1 << 5
	primitive = 0 << 5

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
	typeVideotextString   = 0x15
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

func Parse(out indenter.Indenter, data []byte) {
	for len(data) > 0 {
		data = parseOne(out, data)
	}
}

func parseOne(out indenter.Indenter, data []byte) (rest []byte) {
	if len(data) < 2 {
		panic("short DER read, need at least two bytes, got " + strconv.Itoa(len(data)))
	}

	typeByte := data[0]
	typeTag := typeByte & 0x1F
	typeComposed := typeByte & 0x20
	typeClass := typeByte & 0xC0

	if typeTag == typeIsLongFormTag {
		panic("Long form DER types not implemented")
	}

	switch typeClass {
	case classUniversal:
	case classApplication:
		out.Print("APPLICATION ")
	case classContextSpecific:
		out.Print("CONTEXT-SPECIFIC ")
	case classPrivate:
		out.Print("PRIVATE ")
	}

	switch typeComposed {
	case primitive:
		out.Print("PRIMITIVE ")
	case composed:
		out.Print("COMPOSED ")
	}

	contentLen, rest := decodeLength(data[1:])
	if len(rest) < contentLen {
		panic("Short content, need " + strconv.Itoa(contentLen) + " bytes but have " + strconv.Itoa(len(rest)))
	}
	content := rest[:contentLen]
	rest = rest[contentLen:]

	switch typeByte {
	case typeEndOfContent | primitive:
		if contentLen != 0 {
			panic("End-of-content had unexpected length " + strconv.Itoa(contentLen))
		}
		out.Println("END-OF-CONTENT")
	case typeBoolean | primitive:
		if contentLen != 1 {
			panic("Boolean had unexpected length " + strconv.Itoa(contentLen))
		}
		if content[0] == byte(0) {
			out.Println("BOOLEAN FALSE")
		} else {
			out.Println("BOOLEAN FALSE")
		}
	case typeInteger | primitive:
		if contentLen < 1 {
			panic("Integer had no content")
		}
		if contentLen > 8 {
			panic("Can't handle integer of " + strconv.Itoa(contentLen) + " octets")
		}
		value := int64(0)
		if content[0]&0x80 == 0 {
			// positive number
			for _, v := range content {
				value *= 256
				value += int64(v)
			}
		} else {
			// negative number
			for i, v := range content {
				value *= 256
				if i == 0 && v == 0xff {
					// skip
					continue
				}
				value -= int64(0xff^v) + 1
				// fixme - test with multiple octets
			}
		}
		out.Println("INTEGER", value)
	case typeSequence | composed:
		out.Println("SEQUENCE")
		Parse(out.NextLevel(), content)
	default:
		out.Printf("UNHANDLED CLASS:%x COMPOSED:%x TAG:%02x LENGTH:%d DATA:%s\n",
			typeClass, typeComposed, typeTag, contentLen, base64.StdEncoding.EncodeToString(content))
	}

	return rest
}

func decodeLength(data []byte) (length int, rest []byte) {
	firstByte := data[0]
	if firstByte&0x80 != 0 {
		numToRead := int(firstByte ^ 0x80)
		if len(data)-1 < numToRead {
			panic("Can't satisfy request to read " + strconv.Itoa(numToRead) + " bytes to get length")
		}
		length := 0
		for i := 0; i < numToRead; i++ {
			length *= 256
			length += int(data[1+i])
		}
		return length, data[1+numToRead:]
	} else {
		return int(firstByte), data[1:]
	}
}
