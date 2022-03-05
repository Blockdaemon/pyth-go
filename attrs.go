package pyth

import (
	"bytes"
	"fmt"
	"io"
)

// AttrsMap is a list of string key-value pairs with stable order.
type AttrsMap struct {
	kvps [][2]string
}

// NewAttrsMap returns a new attribute map with an initial arbitrary order.
func NewAttrsMap(fromGo map[string]string) (out AttrsMap, err error) {
	if fromGo != nil {
		for k, v := range fromGo {
			if len(k) > 0xFF {
				return out, fmt.Errorf("key too long (%d > 0xFF): \"%s\"", len(k), k)
			}
			if len(v) > 0xFF {
				return out, fmt.Errorf("value too long (%d > 0xFF): \"%s\"", len(v), v)
			}
			out.kvps = append(out.kvps, [2]string{k, v})
		}
	}
	return
}

// KVs returns the AttrsMap as an unordered Go map.
func (a AttrsMap) KVs() map[string]string {
	m := make(map[string]string, len(a.kvps))
	for _, kv := range a.kvps {
		m[kv[0]] = kv[1]
	}
	return m
}

// UnmarshalBinary unmarshals AttrsMap from its on-chain format.
//
// Will return an error if it fails to consume the entire provided byte slice.
func (a *AttrsMap) UnmarshalBinary(data []byte) (err error) {
	*a, _, err = ReadAttrsMapFromBinary(bytes.NewReader(data))
	return
}

// ReadAttrsMapFromBinary consumes all bytes from a binary reader,
// returning an AttrsMap and the number of bytes read.
func ReadAttrsMapFromBinary(rd *bytes.Reader) (out AttrsMap, n int, err error) {
	for rd.Len() > 0 {
		key, n2, err := readLPString(rd)
		if err != nil {
			return out, n, err
		}
		n += n2
		val, n3, err := readLPString(rd)
		if err != nil {
			return out, n, err
		}
		n += n3
		out.kvps = append(out.kvps, [2]string{key, val})
	}
	return out, n, nil
}

// MarshalBinary marshals AttrsMap to its on-chain format.
func (a AttrsMap) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	for _, kv := range a.kvps {
		if err := writeLPString(&buf, kv[0]); err != nil {
			return nil, err
		}
		if err := writeLPString(&buf, kv[1]); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// readLPString returns a length-prefixed string as seen in AttrsMap.
func readLPString(rd *bytes.Reader) (s string, n int, err error) {
	var strLen byte
	strLen, err = rd.ReadByte()
	if err != nil {
		return
	}
	val := make([]byte, strLen)
	n, err = rd.Read(val)
	n += 1
	s = string(val)
	return
}

// writeLPString writes a length-prefixed string as seen in AttrsMap.
func writeLPString(wr io.Writer, s string) error {
	if len(s) > 0xFF {
		return fmt.Errorf("string too long (%d)", len(s))
	}
	if _, err := wr.Write([]byte{uint8(len(s))}); err != nil {
		return err
	}
	_, err := wr.Write([]byte(s))
	return err
}
