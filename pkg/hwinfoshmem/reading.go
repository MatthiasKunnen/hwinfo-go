package hwinfoshmem

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

	// Original Label in English language.
	OriginalLabelAscii HwinfoSensorStringAscii

	// Displayed label which might have been renamed by the user. Use UserLabel instead.
	UserLabelAscii HwinfoSensorStringAscii

	// The unit of the reading. E.g. °C, RPM. Use Unit instead.
	UnitAscii HwinfoUnitStringAscii

	// The value of the reading.
	Value HwinfoFloat64

	// The minimum value of the reading.
	ValueMin HwinfoFloat64

	// The maximum value of the reading.
	ValueMax HwinfoFloat64

	// The average value of the reading.
	ValueAvg HwinfoFloat64

	// Displayed label which might have been renamed by the user.
	UserLabel HwinfoSensorStringUtf8

	// The unit of the reading. E.g. °C, RPM.
	Unit HwinfoUnitStringUtf8
}
