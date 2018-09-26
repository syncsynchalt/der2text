package pem_test

import (
	"github.com/syncsynchalt/der2text/test"
	. "github.com/syncsynchalt/der2text/write/pem"
	"testing"
)

func b(s string) []byte {
	return []byte(s)
}

func s(b []byte) string {
	return string(b)
}

func TestPemEncodeEmpty(t *testing.T) {
	out := PemEncode("X", b(""))
	test.Equals(t, "-----BEGIN X-----\n-----END X-----\n", s(out))
}

func TestPemEncode(t *testing.T) {
	out := PemEncode("X", b("a"))
	test.Equals(t, "-----BEGIN X-----\nYQ==\n-----END X-----\n", s(out))
	out = PemEncode("X", b("ab"))
	test.Equals(t, "-----BEGIN X-----\nYWI=\n-----END X-----\n", s(out))
	out = PemEncode("X", b("abc"))
	test.Equals(t, "-----BEGIN X-----\nYWJj\n-----END X-----\n", s(out))
}

func TestPemEncodeMultiLabel(t *testing.T) {
	out := PemEncode("X Y Z", b(""))
	test.Equals(t, "-----BEGIN X Y Z-----\n-----END X Y Z-----\n", s(out))
}

func TestPemEncodeEdges(t *testing.T) {
	out := PemEncode("X", b("012345678901234567890123456789012345678901234"))
	test.Equals(t, `-----BEGIN X-----
MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0
-----END X-----
`, s(out))
	out = PemEncode("X", b("0123456789012345678901234567890123456789012345"))
	test.Equals(t, `-----BEGIN X-----
MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NQ==
-----END X-----
`, s(out))
	out = PemEncode("X", b("01234567890123456789012345678901234567890123456"))
	test.Equals(t, `-----BEGIN X-----
MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY=
-----END X-----
`, s(out))
	out = PemEncode("X", b("012345678901234567890123456789012345678901234567"))
	test.Equals(t, `-----BEGIN X-----
MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3
-----END X-----
`, s(out))
	out = PemEncode("X", b("0123456789012345678901234567890123456789012345678"))
	test.Equals(t, `-----BEGIN X-----
MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3
OA==
-----END X-----
`, s(out))
}
