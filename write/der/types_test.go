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
