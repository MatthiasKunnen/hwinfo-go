package hwinfoshmem

// HwinfoSensor can be seen as a way to group readings.
type HwinfoSensor struct {
	// A unique Sensor ID
	SensorId uint32

	// The instance of the sensor (together with SensorId forms a unique ID)
	SensorInstance uint32

	// Original name of sensor in English.
	SensorNameOriginalAscii HwinfoSensorStringAscii

	// Display name of sensor. Might be translated or renamed by user.
	SensorNameAscii HwinfoSensorStringAscii

	// Display name of the sensor. Might be renamed by the user.
	// E.g.
	//   - GIGABYTE B650E AORUS MASTER (ITE IT8689E)
	//   - CPU [#0]: AMD Ryzen 9 7950X
	SensorNameUtf8 HwinfoSensorStringUtf8
}
