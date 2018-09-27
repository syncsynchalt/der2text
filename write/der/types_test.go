package der_test

import (
	"github.com/syncsynchalt/der2text/test"
	. "github.com/syncsynchalt/der2text/write/der"
	"testing"
)

const (
	u = "UNIVERSAL"
	p = "PRIMITIVE"
)

func b(s string) []byte {
	return []byte(s)
}

func TestDerEOC(t *testing.T) {
	data, err := WriteEndOfContent(u, p, "END-OF-CONTENT")
	test.Ok(t, err)
	test.Equals(t, b("\x00\x00"), data)
}

func TestDerBoolean(t *testing.T) {
	data, err := WriteBoolean(u, p, "BOOLEAN", "FALSE")
	test.Ok(t, err)
	test.Equals(t, b("\x01\x01\x00"), data)

	data, err = WriteBoolean(u, p, "BOOLEAN", "TRUE")
	test.Ok(t, err)
	test.Equals(t, b("\x01\x01\x01"), data)

	data, err = WriteBoolean(u, p, "BOOLEAN", "false")
	test.Equals(t, "unrecognized boolean value false", err.Error())
}

func TestDerIntegerAtom(t *testing.T) {
	data, err := WriteInteger(u, p, "INTEGER", "0")
	test.Ok(t, err)
	test.Equals(t, b("\x02\x01\x00"), data)

	data, err = WriteInteger(u, p, "INTEGER", "127")
	test.Ok(t, err)
	test.Equals(t, b("\x02\x01\x7F"), data)

	data, err = WriteInteger(u, p, "INTEGER", "128")
	test.Ok(t, err)
	test.Equals(t, b("\x02\x02\x00\x80"), data)

	data, err = WriteInteger(u, p, "INTEGER", "256")
	test.Ok(t, err)
	test.Equals(t, b("\x02\x02\x01\x00"), data)

	data, err = WriteInteger(u, p, "INTEGER", "-128")
	test.Ok(t, err)
	test.Equals(t, b("\x02\x01\x80"), data)

	data, err = WriteInteger(u, p, "INTEGER", "-129")
	test.Ok(t, err)
	test.Equals(t, b("\x02\x02\xFF\x7F"), data)
}

func TestDerGeneric(t *testing.T) {
	data, err := WriteGeneric(u, p, "INTEGER", "\x01\x02\x03")
	test.Ok(t, err)
	test.Equals(t, b("\x02\x03\x01\x02\x03"), data)
}

func TestDerBitString(t *testing.T) {
	data, err := WriteBitstring(u, p, "BITSTRING", "PAD=0", "\x01\x02\x03")
	test.Ok(t, err)
	test.Equals(t, b("\x03\x04\x00\x01\x02\x03"), data)
}

func TestDerOctetString(t *testing.T) {
	data, err := WriteGeneric(u, p, "OCTETSTRING", "\x01\x02\x03")
	test.Ok(t, err)
	test.Equals(t, b("\x04\x03\x01\x02\x03"), data)
}

func TestDerNullString(t *testing.T) {
	data, err := WriteNull(u, p, "NULL")
	test.Ok(t, err)
	test.Equals(t, b("\x05\x00"), data)
}

func TestDerOid(t *testing.T) {
	data, err := WriteOid(u, p, "OID", "1")
	test.Ok(t, err)
	test.Equals(t, b("\x06\x01\x28"), data)

	data, err = WriteOid(u, p, "OID", "1.2")
	test.Ok(t, err)
	test.Equals(t, b("\x06\x01\x2A"), data)

	data, err = WriteOid(u, p, "OID", "1.2.3")
	test.Ok(t, err)
	test.Equals(t, b("\x06\x02\x2A\x03"), data)

	data, err = WriteOid(u, p, "OID", "1.2.3.4")
	test.Ok(t, err)
	test.Equals(t, b("\x06\x03\x2A\x03\x04"), data)
}

func TestDerOidLong(t *testing.T) {
	data, err := WriteOid(u, p, "OID", "1.2.378")
	test.Ok(t, err)
	test.Equals(t, b("\x06\x03\x2A\x82\x7A"), data)

	data, err = WriteOid(u, p, "OID", "1.2.378.423443")
	test.Ok(t, err)
	test.Equals(t, b("\x06\x06\x2A\x82\x7A\x99\xEC\x13"), data)
}

func TestDerRelativeOid(t *testing.T) {
	data, err := WriteRelativeOid(u, p, "RELATIVEOID", "1.2.378")
	test.Ok(t, err)
	test.Equals(t, b("\x0D\x04\x01\x02\x82\x7A"), data)
}

func TestDerUniversalString(t *testing.T) {
	data, err := WriteUniversalString(u, p, "UNIVERSALSTRING", "foo")
	test.Ok(t, err)
	test.Equals(t, b("\x1C\x0C\x00\x00\x00f\x00\x00\x00o\x00\x00\x00o"), data)

	data, err = WriteUniversalString(u, p, "UNIVERSALSTRING", "😏")
	test.Ok(t, err)
	test.Equals(t, b("\x1C\x04\x00\x01\xF6\x0F"), data)
}

func TestDerBMPString(t *testing.T) {
	data, err := WriteBMPString(u, p, "BMPSTRING", "foo")
	test.Ok(t, err)
	test.Equals(t, b("\x1E\x06\x00f\x00o\x00o"), data)

	data, err = WriteBMPString(u, p, "BMPSTRING", "😏")
	test.Ok(t, err)
	test.Equals(t, b("\x1E\x04\xD8\x3D\xDE\x0F"), data)
}
