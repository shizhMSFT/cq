package cq

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/shizhMSFT/cq/internal/cbor"
)

func Print(r io.Reader) error {
	return print(0, r, 0, "")
}

func PrintBytes(c []byte) error {
	return Print(bytes.NewReader(c))
}

func print(indent int, r io.Reader, width int, prefix string) error {
	majorType, count, content, err := cbor.ReadHeader(r)
	if err != nil {
		return err
	}

	switch majorType {
	case 0, 1:
		n := int64(count)
		if majorType == 1 {
			n = -1 - n
		}
		desc := fmt.Sprintf("Integer: %d", n)
		println(indent, content, width, prefix+desc)
	case 2:
		desc := fmt.Sprintf("Binary string: %d bytes", count)
		println(indent, content, width, prefix+desc)
		return printString(indent+1, r, count)
	case 3:
		desc := fmt.Sprintf("UTF-8 text: %d bytes", count)
		println(indent, content, width, prefix+desc)
		return printString(indent+1, r, count)
	case 4:
		desc := fmt.Sprintf("Array of length %d", count)
		println(indent, content, width, prefix+desc)
		for i := uint64(0); i < count; i++ {
			if err := print(indent+1, r, 0, ""); err != nil {
				return err
			}
		}
	case 5:
		desc := fmt.Sprintf("Map of size %d", count)
		println(indent, content, 0, prefix+desc)
		for i := uint64(0); i < count; i++ {
			if err := print(indent+1, r, 3, "Key:   "); err != nil {
				return err
			}
			if err := print(indent+1, r, 3, "Value: "); err != nil {
				return err
			}
		}
	case 6:
		return printTag(indent, r, content, count, prefix)
	default:
		return fmt.Errorf("major type %d is not supported", majorType)
	}
	return nil
}

func printString(indent int, r io.Reader, count uint64) error {
	var buf [16]byte
	for count > 0 {
		n := count
		if n > 16 {
			n = 16
		}
		line := buf[:n]
		if _, err := io.ReadFull(r, line); err != nil {
			return err
		}
		b := strings.Builder{}
		for _, c := range line {
			if c > 0x1f && c < 0x7f {
				b.WriteByte(c)
			} else {
				b.WriteByte('.')
			}
		}
		println(indent, line, 16, b.String())
		count -= n
	}
	return nil
}

func printTag(indent int, r io.Reader, content []byte, tag uint64, prefix string) error {
	desc := fmt.Sprintf("%sTag %d", prefix, tag)
	switch tag {
	case 1:
		desc += ": datetime"
		println(indent, content, 0, desc)
		return printDateTime(indent, r)
	case 18:
		desc += ": cose-sign1"
	}
	println(indent, content, 0, desc)
	return print(indent, r, 0, "")
}

func printDateTime(indent int, r io.Reader) error {
	majorType, count, content, err := cbor.ReadHeader(r)
	if err != nil {
		return err
	}
	if majorType != 0 {
		return fmt.Errorf("datetime type %d not supported", majorType)
	}

	t := time.Unix(int64(count), 0).UTC()
	desc := fmt.Sprintf("UNIX epoch: %d -> %s", count, t.Format(time.RFC3339))
	println(indent, content, 0, desc)
	return nil
}

func println(indent int, content []byte, width int, description string) {
	for i := 0; i < indent; i++ {
		fmt.Printf("   ")
	}
	for _, c := range content {
		fmt.Printf("%02x ", c)
	}
	padding := width - len(content)
	for i := 0; i < padding; i++ {
		fmt.Printf("   ")
	}
	fmt.Println("--", description)
}
