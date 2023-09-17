package util

import "bytes"

// NulTerminatedUtf8ByteArrayToString converts a byte slice containing UTF-8 bytes to a string.
// The resulting string will consist of all characters until a nul byte is encountered or the end of
// the slice is reached.
func NulTerminatedUtf8ByteArrayToString(data []byte) string {
	nulByteIndex := bytes.IndexByte(data, 0)
	if nulByteIndex < 0 {
		return string(data[:])
	} else {
		return string(data[:nulByteIndex])
	}
}
