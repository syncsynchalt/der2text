package der

import (
	"github.com/syncsynchalt/der2text/test"
	"testing"
)

func TestClassMask(t *testing.T) {
	test.Equals(t, true, IsValidClass("UNIVERSAL"))
	test.Equals(t, false, IsValidClass("Universal"))
	test.Equals(t, false, IsValidClass(""))
	test.Equals(t, true, IsValidClass("APPLICATION"))
	test.Equals(t, true, IsValidClass("CONTEXT-SPECIFIC"))
	test.Equals(t, true, IsValidClass("PRIVATE"))
}

func TestConstructedMask(t *testing.T) {
	test.Equals(t, false, IsValidConstructed(""))
	test.Equals(t, false, IsValidConstructed("xxx"))
	test.Equals(t, false, IsValidConstructed("primitive"))
	test.Equals(t, true, IsValidConstructed("PRIMITIVE"))
	test.Equals(t, true, IsValidConstructed("CONSTRUCTED"))
}

func TestTypeValue(t *testing.T) {
	universalTypes := []string{"END-OF-CONTENT", "BOOLEAN", "INTEGER", "BITSTRING", "OCTETSTRING", "NULL", "OID",
		"OBJECTDESCRIPTION", "EXTERNAL", "REAL", "ENUMERATED", "EMBEDDED-PDV", "UTF8STRING", "RELATIVEOID",
		"", "", "SEQUENCE", "SET", "NUMERICSTRING", "PRINTABLESTRING", "T61STRING", "VIDEOTEXSTRING",
		"IA5STRING", "UTCTIME", "GENERALIZEDTIME", "GRAPHICSTRING", "VISIBLESTRING", "GENERALSTRING",
		"UNIVERSALSTRING", "CHARACTERSTRING", "BMPSTRING"}
	for i, v := range universalTypes {
		t.Log("Checking", v)
		if len(v) == 0 {
			continue
		}
		bytes, err := typeBytes(v)
		test.Ok(t, err)
		test.Equals(t, []byte{byte(i)}, bytes)
	}
}

func TestTypeValueUnhandled(t *testing.T) {
	bytes, err := typeBytes("UNHANDLED-TAG=23")
	test.Ok(t, err)
	test.Equals(t, []byte{23}, bytes)

	bytes, err = typeBytes("UNHANDLED-TAG=12x")
	test.Equals(t, "strconv.Atoi: parsing \"12x\": invalid syntax", err.Error())

	bytes, err = typeBytes("UNHANDLED-TAG=16384")
	test.Equals(t, "Tag number of 16384 not implemented (max of 0x3fff)", err.Error())

	bytes, err = typeBytes("UNHANDLED-TAG=123")
	test.Ok(t, err)
	test.Equals(t, []byte{0x1f, 123}, bytes)

	bytes, err = typeBytes("UNHANDLED-TAG=1234")
	test.Ok(t, err)
	test.Equals(t, []byte{0x1f, 0x89, 0x52}, bytes)
}

func TestLength(t *testing.T) {
	bytes := lengthBytes(0)
	test.Equals(t, []byte{0x00}, bytes)

	bytes = lengthBytes(1)
	test.Equals(t, []byte{0x01}, bytes)

	bytes = lengthBytes(127)
	test.Equals(t, []byte{0x7f}, bytes)

	bytes = lengthBytes(128)
	test.Equals(t, []byte{0x81, 0x80}, bytes)

	bytes = lengthBytes(5000)
	test.Equals(t, []byte{0x82, 0x13, 0x88}, bytes)
}
