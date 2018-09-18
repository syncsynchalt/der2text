package pem

import (
	"encoding/base64"
	"strings"

	"der2text/der"
	"der2text/indenter"
)

func Parse(out indenter.Indenter, data []byte) {
	str := string(data)

	str = strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' {
			return -1
		} else {
			return r
		}
	}, str)

	if len(str) < 11 || str[:11] != "-----BEGIN " {
		panic("Unable to parse PEM header")
	}

	eot := strings.IndexRune(str[11:], '-')
	typ := str[11 : 11+eot]

	head := "-----BEGIN " + typ + "-----"
	tail := "-----END " + typ + "-----"
	if !strings.HasPrefix(str, head) {
		panic("PEM doesn't have expected prefix " + head)
	}
	if !strings.HasSuffix(str, tail) {
		panic("PEM doesn't have expected suffix " + tail)
	}

	out.Println("PEM ENCODED", typ)

	b64 := str[len(head) : len(str)-len(tail)-1]
	for len(b64)%4 != 0 {
		b64 += "="
	}
	derData, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	der.Parse(out.NextLevel(), derData)
}
