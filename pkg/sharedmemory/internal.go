package sharedmemory

import "unsafe"

func GetHwinfoHeaderFromStartPointer(pointer uintptr) *HwinfoHeader {
	return (*HwinfoHeader)(unsafe.Pointer(pointer))
}

func GetSensorsFromStartPointer(info *HwinfoHeader, pointer uintptr) []*HwinfoSensor {
	sensors := make([]*HwinfoSensor, info.SensorAmount)

	for i := uint32(0); i < info.SensorAmount; i++ {
		offset := pointer + uintptr(info.SensorSectionOffset) + uintptr(i)*uintptr(info.SensorSize)
		sensorElement := (*HwinfoSensor)(unsafe.Pointer(offset))
		sensors[i] = sensorElement
	}

	return sensors
}

func GetReadingsFromStartPointer(info *HwinfoHeader, pointer uintptr) []*HwinfoReading {
	readings := make([]*HwinfoReading, info.ReadingAmount)

	for i := uint32(0); i < info.ReadingAmount; i++ {
		offset := pointer + uintptr(info.ReadingSectionOffset) + uintptr(i)*uintptr(info.ReadingSize)
		readingElement := (*HwinfoReading)(unsafe.Pointer(offset))
		readings[i] = readingElement
	}

	return readings
}

func GetReadingsById(info *HwinfoHeader, pointer uintptr, readingIds []ReadingIdSensorCombo) []*HwinfoReading {
	readings := make([]*HwinfoReading, 0)

	for i := uint32(0); i < info.ReadingAmount; i++ {
		offset := pointer + uintptr(info.ReadingSectionOffset) + uintptr(i)*uintptr(info.ReadingSize)
		sensorIndex := *(*uint32)(unsafe.Pointer(offset + uintptr(4)))
		readingId := *(*uint32)(unsafe.Pointer(offset + uintptr(8)))

		for _, indexInfo := range readingIds {
			if indexInfo.SensorIndex == sensorIndex && indexInfo.Id == readingId {
				reading := (*HwinfoReading)(unsafe.Pointer(offset))
				readings = append(readings, reading)
			}
		}
	}

	return readings
}
