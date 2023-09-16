package sharedmemory

import "unsafe"

type Reader struct {
	GetPointer func() (uintptr, error)
}

// GetHeader returns the header of the shared memory. Make sure to lock using Lock().
func (reader *Reader) GetHeader() (*HwinfoHeader, error) {
	pointer, err := reader.GetPointer()
	if err != nil {
		return nil, err
	}

	return (*HwinfoHeader)(unsafe.Pointer(pointer)), nil
}

// GetSensors returns the sensors that are reported by HWiNFO.
// Make sure that the given HwinfoHeader is current, meaning that the lock was held when calling
// GetHeader and stays held while calling this function and processing its results.
func (reader *Reader) GetSensors(info *HwinfoHeader) ([]*HwinfoSensor, error) {
	pointer, err := reader.GetPointer()
	if err != nil {
		return nil, err
	}

	sensors := make([]*HwinfoSensor, info.SensorAmount)

	for i := uint32(0); i < info.SensorAmount; i++ {
		offset := pointer + uintptr(info.SensorSectionOffset) + uintptr(i)*uintptr(info.SensorSize)
		sensorElement := (*HwinfoSensor)(unsafe.Pointer(offset))
		sensors[i] = sensorElement
	}

	return sensors, nil
}

// GetReadings returns all HWiNFO readings.
func (reader *Reader) GetReadings(info *HwinfoHeader) ([]*HwinfoReading, error) {
	pointer, err := reader.GetPointer()
	if err != nil {
		return nil, err
	}

	readings := make([]*HwinfoReading, info.ReadingAmount)

	for i := uint32(0); i < info.ReadingAmount; i++ {
		offset := pointer + uintptr(info.ReadingSectionOffset) + uintptr(i)*uintptr(info.ReadingSize)
		readingElement := (*HwinfoReading)(unsafe.Pointer(offset))
		readings[i] = readingElement
	}

	return readings, nil
}

type ReadingIdSensorCombo struct {
	Id          uint32
	SensorIndex uint32
}

// GetReadingsById returns the readings that match the given sensor index/id combinations.
func (reader *Reader) GetReadingsById(info *HwinfoHeader, readingIds []ReadingIdSensorCombo) ([]*HwinfoReading, error) {
	readings := make([]*HwinfoReading, 0)
	pointer, err := reader.GetPointer()
	if err != nil {
		return nil, err
	}

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

	return readings, nil
}
