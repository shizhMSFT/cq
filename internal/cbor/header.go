package cbor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// ReadHeader reads the header of a CBOR data item.
// It returns the major type, the count, and the content of the header.
func ReadHeader(r io.Reader) (byte, uint64, []byte, error) {
	contentBuffer := bytes.NewBuffer(nil)
	r = io.TeeReader(r, contentBuffer)

	var header [1]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return 0, 0, nil, err
	}

	majorType := header[0] >> 5
	count := uint64(header[0] & 0x1f)
	if count > 27 {
		return 0, 0, nil, fmt.Errorf("invalid count: %d", count)
	}
	switch count {
	case 24:
		var counts [1]byte
		if _, err := io.ReadFull(r, counts[:]); err != nil {
			return 0, 0, nil, err
		}
		count = uint64(counts[0])
	case 25:
		var counts [2]byte
		if _, err := io.ReadFull(r, counts[:]); err != nil {
			return 0, 0, nil, err
		}
		count = uint64(binary.BigEndian.Uint16(counts[:]))
	case 26:
		var counts [4]byte
		if _, err := io.ReadFull(r, counts[:]); err != nil {
			return 0, 0, nil, err
		}
		count = uint64(binary.BigEndian.Uint32(counts[:]))
	case 27:
		var counts [8]byte
		if _, err := io.ReadFull(r, counts[:]); err != nil {
			return 0, 0, nil, err
		}
		count = binary.BigEndian.Uint64(counts[:])
	}

	return majorType, count, contentBuffer.Bytes(), nil
}
