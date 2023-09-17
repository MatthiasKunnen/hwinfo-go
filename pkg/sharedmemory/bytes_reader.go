package sharedmemory

import "unsafe"

// BytesReader allows for extracting the header, sensors, and readings from a copy of HWiNFO's
// shared memory.
// A copy can be made using [MemoryReader.Copy].
type BytesReader struct {
	Bytes []byte
	Data  Reader
}

func NewBytesReader(bytes []byte) *BytesReader {
	bytesReader := &BytesReader{
		Bytes: bytes,
	}
	bytesReader.Data.GetPointer = func() (uintptr, error) {
		return uintptr(unsafe.Pointer(&bytesReader.Bytes[0])), nil
	}

	return bytesReader
}
