package der

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"strconv"
	"strings"
)

func WriteEndOfContent(class, constructed, typ string) ([]byte, error) {
	return makeDer(class, constructed, typ, []byte{})
}

func WriteBoolean(class, constructed, typ, value string) ([]byte, error) {
	switch value {
	case "TRUE":
		return makeDer(class, constructed, typ, []byte{0x01})
	case "FALSE":
		return makeDer(class, constructed, typ, []byte{0x00})
	default:
		return nil, fmt.Errorf("unrecognized boolean value %s", value)
	}
}

func WriteInteger(class, constructed, typ, atom string) ([]byte, error) {

	asInt, err := strconv.ParseInt(atom, 10, 64)
	if err != nil {
		return nil, err
	}

	// convert int to bytes
	asUint := uint64(asInt)
	payload := make([]byte, 64/8)
	binary.BigEndian.PutUint64(payload, asUint)

	if asInt == 0 {
		payload = []byte{0x00}
	} else if asInt > 0 {
		for len(payload) > 1 && payload[0] == 0x00 && payload[1]&0x80 == 0 {
			payload = payload[1:]
		}
	} else if asInt < 0 {
		for len(payload) > 1 && payload[0] == 0xFF && payload[1]&0x80 == 0x80 {
			payload = payload[1:]
		}
	}

	return makeDer(class, constructed, typ, payload)
}

func WriteGeneric(class, constructed, typ, payload string) ([]byte, error) {
	return makeDer(class, constructed, typ, []byte(payload))
}

func WriteBitstring(class, constructed, typ, padding, payload string) ([]byte, error) {
	if len(padding) <= 4 || padding[:4] != "PAD=" {
		return nil, fmt.Errorf("Unrecognized padding %s (expected PAD=nnn)", padding)
	}
	padAsNum, err := strconv.Atoi(padding[4:])
	if err != nil {
		return nil, fmt.Errorf("Unrecognized number %s (expected PAD=nnn)", padding)
	}
	realPayload := make([]byte, 0)
	realPayload = append(realPayload, byte(padAsNum))
	realPayload = append(realPayload, payload...)
	return makeDer(class, constructed, typ, realPayload)
}

func WriteNull(class, constructed, typ string) ([]byte, error) {
	return makeDer(class, constructed, typ, []byte{})
}

func encodeOidParts(parts []string) ([]byte, error) {
	data := make([]byte, 0)
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		reversed := make([]byte, 0)
		// convert to little-endian bytes, 7 bits per byte with 8th bit set
		for num != 0 {
			reversed = append(reversed, byte(0x80|(num&0x7f)))
			num = num >> 7
		}
		// lay it in data backwards to big-endian it
		for len(reversed) != 0 {
			data = append(data, reversed[len(reversed)-1])
			reversed = reversed[:len(reversed)-1]
		}
		// drop the 8th bit on the last byte
		data[len(data)-1] = data[len(data)-1] & 0x7f
	}
	return data, nil
}

func WriteOid(class, constructed, typ, oid string) ([]byte, error) {
	parts := strings.Split(oid, ".")

	var b1, b2 int
	var err error
	if len(parts) > 0 {
		b1, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
	}
	if len(parts) > 1 {
		b2, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
	}
	payload := make([]byte, 0)
	payload = append(payload, byte(40*b1+b2))
	if len(parts) >= 2 {
		d, err := encodeOidParts(parts[2:])
		if err != nil {
			return nil, err
		}
		payload = append(payload, d...)
	}

	return makeDer(class, constructed, typ, payload)
}

func WriteRelativeOid(class, constructed, typ, oid string) ([]byte, error) {
	parts := strings.Split(oid, ".")
	payload, err := encodeOidParts(parts)
	if err != nil {
		return nil, err
	}
	return makeDer(class, constructed, typ, payload)
}

// aka UTF-32BE / UCS-4 / etc
func WriteUniversalString(class, constructed, typ, str string) ([]byte, error) {
	encoder := utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewEncoder()
	payload, err := encoder.Bytes([]byte(str))
	if err != nil {
		return nil, err
	}
	return makeDer(class, constructed, typ, payload)
}

// aka UTF-16BE / UCS-2 / etc
func WriteBMPString(class, constructed, typ, str string) ([]byte, error) {
	encoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()
	payload, err := encoder.Bytes([]byte(str))
	if err != nil {
		return nil, err
	}
	return makeDer(class, constructed, typ, payload)
}
