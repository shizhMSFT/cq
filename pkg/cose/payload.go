package cose

import (
	"errors"
	"fmt"
	"io"

	"github.com/shizhMSFT/cq/internal/cbor"
)

// ExtractPayload reads the payload of a COSE_Sign1 object from r, and stops at
// the end of the payload.
func ExtractPayload(r io.Reader) ([]byte, error) {
	// parse tag
	majorType, count, _, err := cbor.ReadHeader(r)
	if err != nil {
		return nil, err
	}
	if majorType == 6 {
		if count != 18 {
			return nil, errors.New("not a cose-sign1 object")
		}
		majorType, count, _, err = cbor.ReadHeader(r)
		if err != nil {
			return nil, err
		}
	}
	if majorType != 4 && count != 4 {
		return nil, errors.New("invalid cose-sign1 object")
	}

	// skip headers
	for i := 0; i < 2; i++ {
		if err := cbor.Discard(r); err != nil {
			return nil, err
		}
	}

	// extract payload
	majorType, count, _, err = cbor.ReadHeader(r)
	if err != nil {
		return nil, err
	}
	if majorType != 2 {
		return nil, errors.New("invalid cose-sign1 object: payload is not a binary string")
	}
	payload := make([]byte, count)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}

	return payload, nil
}
