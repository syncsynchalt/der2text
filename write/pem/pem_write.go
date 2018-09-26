package pem

import (
	"encoding/base64"
)

func PemEncode(label string, data []byte) []byte {
	var result []byte

	result = append(result, "-----BEGIN "...)
	result = append(result, label...)
	result = append(result, "-----\n"...)

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)
	for len(encoded) > 64 {
		result = append(result, encoded[:64]...)
		result = append(result, '\n')
		encoded = encoded[64:]
	}
	if len(encoded) > 0 {
		result = append(result, encoded...)
		result = append(result, '\n')
	}

	result = append(result, "-----END "...)
	result = append(result, label...)
	result = append(result, "-----\n"...)

	return result
}
