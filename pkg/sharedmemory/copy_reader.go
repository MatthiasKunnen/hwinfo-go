package sharedmemory

import "unsafe"

// CopyReader allows for extracting the header, sensors, and readings from a copy of HWiNFO's shared
// memory. A copy can be made using [MemoryReader.Copy].
type CopyReader struct {
	Bytes []byte
}

func (reader *CopyReader) GetHeader() *HwinfoHeader {
	return GetHwinfoHeaderFromStartPointer(uintptr(unsafe.Pointer(&reader.Bytes[0])))
}

func (reader *CopyReader) GetSensors(info *HwinfoHeader) []*HwinfoSensor {
	return GetSensorsFromStartPointer(info, uintptr(unsafe.Pointer(&reader.Bytes[0])))
}

func (reader *CopyReader) GetReadings(info *HwinfoHeader) []*HwinfoReading {
	return GetReadingsFromStartPointer(info, uintptr(unsafe.Pointer(&reader.Bytes[0])))
}
