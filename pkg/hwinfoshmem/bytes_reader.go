package hwinfoshmem

import "unsafe"

// BytesReader allows for extracting the header, sensors, and readings from a copy of HWiNFO's
// shared memory.
// A copy can be made using [MemoryReader.Copy].
//
// BytesReader has an initializer function, [NewBytesReader].
type BytesReader struct {
	Bytes []byte
	*Reader
}

func NewBytesReader(bytes []byte) *BytesReader {
	bytesReader := &BytesReader{
		Bytes: bytes,
	}
	bytesReader.Reader = &Reader{
		GetPointer: func() (uintptr, error) {
			return uintptr(unsafe.Pointer(&bytesReader.Bytes[0])), nil
		},
	}

	return bytesReader
}
