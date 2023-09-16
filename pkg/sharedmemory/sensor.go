package sharedmemory

// HwinfoSensor can be seen as a way to group readings.
type HwinfoSensor struct {
	// A unique Sensor ID
	SensorId uint32

	// The instance of the sensor (together with SensorId forms a unique ID)
	SensorInstance uint32

	// Original name of sensor in English. ANSI string.
	SensorNameOriginalAnsi [hwinfoSensorStringLength]byte

	// Display name of sensor. Might be translated or renamed by user. ANSI string.
	SensorNameAnsi [hwinfoSensorStringLength]byte

	// Display name of the sensor. Might be renamed by the user.
	// Use GetSensorName to get a string value.
	// UTF-8 string.
	SensorNameUtf8 [hwinfoSensorStringLength]byte
}

// GetSensorName returns the sensor's name as string. E.g.
//   - GIGABYTE B650E AORUS MASTER (ITE IT8689E)
//   - CPU [#0]: AMD Ryzen 9 7950X
func (sensor *HwinfoSensor) GetSensorName() string {
	return string(sensor.SensorNameUtf8[:])
}
