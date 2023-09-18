package util

import "bytes"

// Utf8BytesToString converts UTF-8 bytes to a string stopping at the first nul byte.
func Utf8BytesToString(data []byte) string {
	nulByteIndex := bytes.IndexByte(data, 0)
	if nulByteIndex < 0 {
		return string(data[:])
	} else {
		return string(data[:nulByteIndex])
	}
}
