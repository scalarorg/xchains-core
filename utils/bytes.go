package utils

import (
	"encoding/binary"
	"errors"
)

func BytesToInt64BigEndian(b []byte) (int64, error) {
	if len(b) < 8 {
		return 0, errors.New("byte slice must be at least 8 bytes long")
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}
