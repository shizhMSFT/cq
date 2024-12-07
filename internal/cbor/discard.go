package cbor

import (
	"fmt"
	"io"
)

// Discard reads and discards the next CBOR item from r.
func Discard(r io.Reader) error {
	majorType, count, _, err := ReadHeader(r)
	if err != nil {
		return err
	}

	switch majorType {
	case 0, 1: // integer
	case 2, 3: // binary/text string
		_, err := io.CopyN(io.Discard, r, int64(count))
		return err
	case 4: // array
		for i := uint64(0); i < count; i++ {
			if err := Discard(r); err != nil {
				return err
			}
		}
	case 5: // map
		for i := uint64(0); i < count; i++ {
			if err := Discard(r); err != nil { // key
				return err
			}
			if err := Discard(r); err != nil { // value
				return err
			}
		}
	case 6:
		return discardTag(r, count)
	default:
		return fmt.Errorf("major type %d is not supported", majorType)
	}
	return nil
}

func discardTag(r io.Reader, tag uint64) error {
	switch tag {
	case 1:
		return discardDateTime(r)
	}
	return Discard(r)
}

func discardDateTime(r io.Reader) error {
	majorType, _, _, err := ReadHeader(r)
	if err != nil {
		return err
	}
	if majorType != 0 {
		return fmt.Errorf("datetime type %d not supported", majorType)
	}
	return nil
}
