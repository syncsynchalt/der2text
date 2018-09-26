package der

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

func classMask(class string) (byte, error) {
	switch class {
	case "UNIVERSAL":
		return 0x00, nil
	case "APPLICATION":
		return 0x40, nil
	case "CONTEXT-SPECIFIC":
		return 0x80, nil
	case "PRIVATE":
		return 0xC0, nil
	default:
		return 0x00, errors.New("Unable to find class of type " + class)
	}
}

func IsValidClass(class string) bool {
	_, err := classMask(class)
	return err == nil
}

func constructedMask(constructed string) (byte, error) {
	switch constructed {
	case "PRIMITIVE":
		return 0x00, nil
	case "CONSTRUCTED":
		return 0x20, nil
	default:
		return 0x00, errors.New("Unable to parse construction flag of " + constructed)
	}
}

func IsValidConstructed(constructed string) bool {
	_, err := constructedMask(constructed)
	return err == nil
}

func typeBytes(typ string) ([]byte, error) {
	out := make([]byte, 0, 20)
	switch typ {
	case "END-OF-CONTENT":
		return append(out, 0x00), nil
	case "BOOLEAN":
		return append(out, 0x01), nil
	case "INTEGER":
		return append(out, 0x02), nil
	case "BITSTRING":
		return append(out, 0x03), nil
	case "OCTETSTRING":
		return append(out, 0x04), nil
	case "NULL":
		return append(out, 0x05), nil
	case "OID":
		return append(out, 0x06), nil
	case "OBJECTDESCRIPTION":
		return append(out, 0x07), nil
	case "EXTERNAL":
		return append(out, 0x08), nil
	case "REAL":
		return append(out, 0x09), nil
	case "ENUMERATED":
		return append(out, 0x0A), nil
	case "EMBEDDED-PDV":
		return append(out, 0x0B), nil
	case "UTF8STRING":
		return append(out, 0x0C), nil
	case "RELATIVEOID":
		return append(out, 0x0D), nil
	case "SEQUENCE":
		return append(out, 0x10), nil
	case "SET":
		return append(out, 0x11), nil
	case "NUMERICSTRING":
		return append(out, 0x12), nil
	case "PRINTABLESTRING":
		return append(out, 0x13), nil
	case "T61STRING":
		return append(out, 0x14), nil
	case "VIDEOTEXSTRING":
		return append(out, 0x15), nil
	case "IA5STRING":
		return append(out, 0x16), nil
	case "UTCTIME":
		return append(out, 0x17), nil
	case "GENERALIZEDTIME":
		return append(out, 0x18), nil
	case "GRAPHICSTRING":
		return append(out, 0x19), nil
	case "VISIBLESTRING":
		return append(out, 0x1A), nil
	case "GENERALSTRING":
		return append(out, 0x1B), nil
	case "UNIVERSALSTRING":
		return append(out, 0x1C), nil
	case "CHARACTERSTRING":
		return append(out, 0x1D), nil
	case "BMPSTRING":
		return append(out, 0x1E), nil
	}
	if len(typ) > 14 && typ[:14] == "UNHANDLED-TAG=" {
		tagNum, err := strconv.Atoi(typ[14:])
		if err != nil {
			return nil, err
		}
		if tagNum < 0x1f {
			return append(out, byte(tagNum)), nil
		} else if tagNum <= 0x7f {
			return append(out, 0x1F, 0x00|byte(tagNum)), nil
		} else if tagNum <= 0x3fff {
			return append(out, 0x1F, 0x80|byte(tagNum>>7&0x7F), 0x00|byte(tagNum&0x7F)), nil
		} else {
			return nil, errors.New(fmt.Sprintf("Tag number of %d not implemented (max of 0x3fff)", tagNum))
		}
	}
	return nil, errors.New(fmt.Sprintf("Unrecognized tag %s", typ))
}

func IsValidType(typ string) bool {
	bytes, _ := typeBytes(typ)
	return bytes != nil
}

func lengthBytes(length uint32) []byte {
	if length <= 127 {
		return []byte{byte(length)}
	}
	asBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(asBytes, length)

	// subtract any empty bytes
	numEmpty := 0
	for asBytes[numEmpty] == 0 {
		numEmpty++
	}

	output := make([]byte, 0)
	output = append(output, byte(0x80|(4-numEmpty)))
	for i := numEmpty; i < len(asBytes); i++ {
		output = append(output, asBytes[i])
	}
	return output
}

func makeDer(class, constructed, typ string, payload []byte) ([]byte, error) {
	// type
	bytes, err := typeBytes(typ)
	if err != nil {
		return nil, err
	}
	m1, err := classMask(class)
	if err != nil {
		return nil, err
	}
	m2, err := constructedMask(constructed)
	if err != nil {
		return nil, err
	}
	// first byte of type also has the type tags
	bytes[0] = bytes[0] | m1 | m2

	payloadLen := uint32(len(payload))
	bytes = append(bytes, lengthBytes(payloadLen)...)
	bytes = append(bytes, payload...)
	return bytes, nil
}
