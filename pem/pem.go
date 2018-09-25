package pem

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/syncsynchalt/der2text/der"
	"github.com/syncsynchalt/der2text/indenter"
)

func Parse(out *indenter.Indenter, data []byte) error {
	str := string(data)
	str = strings.Trim(str, " \t\r\n")
	if len(str) < 11 || str[:11] != "-----BEGIN " {
		return errors.New("Unable to parse PEM header")
	}

	eot := strings.IndexRune(str[11:], '-')
	typ := str[11 : 11+eot]

	head := "-----BEGIN " + typ + "-----"
	tail := "-----END " + typ + "-----"
	if !strings.HasPrefix(str, head) {
		return errors.New("PEM doesn't have expected prefix " + head)
	}
	if !strings.HasSuffix(str, tail) {
		return errors.New("PEM doesn't have expected suffix " + tail)
	}

	out.Println("PEM ENCODED", typ)

	b64 := str[len(head) : len(str)-len(tail)]
	b64 = stripSpaces(b64)
	for len(b64)%4 != 0 {
		b64 += "="
	}
	derData, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}

	return der.Parse(out.NextLevel(), derData)
}

func stripSpaces(s string) string {
	return strings.Join(strings.Fields(s), "")
}
