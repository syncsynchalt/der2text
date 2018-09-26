package der

import (
	"encoding/binary"
	"strconv"
)

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
		for len(payload) > 0 && payload[0] == 0x00 && payload[1]&0x80 == 0 {
			payload = payload[1:]
		}
	} else if asInt < 0 {
		for len(payload) > 0 && payload[0] == 0xFF && payload[1]&0x80 == 0x80 {
			payload = payload[1:]
		}
	}

	return makeDer(class, constructed, typ, payload)
}

func WriteIntegerPreserved(class, constructed, typ, payload string) ([]byte, error) {
	return makeDer(class, constructed, typ, []byte(payload))
}
