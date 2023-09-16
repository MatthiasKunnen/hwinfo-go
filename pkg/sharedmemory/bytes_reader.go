package sharedmemory

import "unsafe"

// BytesReader allows for extracting the header, sensors, and readings from a copy of HWiNFO's
// shared memory.
// A copy can be made using [MemoryReader.Copy].
type BytesReader struct {
	Bytes []byte
}

func (reader *BytesReader) GetHeader() *HwinfoHeader {
	return GetHwinfoHeaderFromStartPointer(uintptr(unsafe.Pointer(&reader.Bytes[0])))
}

func (reader *BytesReader) GetSensors(info *HwinfoHeader) []*HwinfoSensor {
	return GetSensorsFromStartPointer(info, uintptr(unsafe.Pointer(&reader.Bytes[0])))
}

func (reader *BytesReader) GetReadings(info *HwinfoHeader) []*HwinfoReading {
	return GetReadingsFromStartPointer(info, uintptr(unsafe.Pointer(&reader.Bytes[0])))
}
