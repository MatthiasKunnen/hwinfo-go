package sharedmemory

import (
	"encoding/binary"
	"github.com/MatthiasKunnen/hwinfo-go/pkg/util"
	"math"
)

type ReadingType uint32

const (
	SENSOR_TYPE_NONE ReadingType = iota
	SENSOR_TYPE_TEMP
	SENSOR_TYPE_VOLT
	SENSOR_TYPE_FAN
	SENSOR_TYPE_CURRENT
	SENSOR_TYPE_POWER
	SENSOR_TYPE_CLOCK
	SENSOR_TYPE_USAGE
	SENSOR_TYPE_OTHER
)

type HwinfoReading struct {
	// The type of reading.
	Type ReadingType

	// Index of the Sensor this reading belongs to.
	SensorIndex uint32

	// A unique ID of the reading within a particular sensor.
	Id uint32

	// Original Label in English language (ANSI).
	OriginalLabelAnsi [hwinfoSensorStringLength]byte

	// Displayed label which might have been renamed by the user (ANSI). Use UserLabelUtf8 instead.
	UserLabelAnsi [hwinfoSensorStringLength]byte

	// The unit of the reading. E.g. °C, RPM (ANSI). Use UnitUtf8 instead.
	UnitAnsi [hwinfoUnitStringLength]byte

	// The value of the reading. Get the actual value, instead of the byte array, using GetValue.
	Value [8]byte

	// The minimum value of the reading. Get the actual value, instead of the byte array, using GetValueMin.
	ValueMin [8]byte

	// The maximum value of the reading. Get the actual value, instead of the byte array, using GetValueMax.
	ValueMax [8]byte

	// The average value of the reading. Get the actual value, instead of the byte array, using GetValueAvg.
	ValueAvg [8]byte

	// Displayed label which might have been renamed by the user (UTF-8). Get the label name using GetUserLabel.
	UserLabelUtf8 [hwinfoSensorStringLength]byte

	// The unit of the reading. E.g. °C, RPM (UTF-8). Get a string using GetUnit.
	UnitUtf8 [hwinfoUnitStringLength]byte
}

func (reading *HwinfoReading) GetUserLabel() string {
	return util.NulTerminatedUtf8ByteArrayToString(reading.UserLabelUtf8[:])
}

func (reading *HwinfoReading) GetUnit() string {
	return util.NulTerminatedUtf8ByteArrayToString(reading.UnitUtf8[:])
}

func (reading *HwinfoReading) GetValue() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(reading.Value[:]))
}

func (reading *HwinfoReading) GetValueMin() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(reading.ValueMin[:]))
}

func (reading *HwinfoReading) GetValueMax() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(reading.ValueMax[:]))
}

func (reading *HwinfoReading) GetValueAvg() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(reading.ValueAvg[:]))
}
