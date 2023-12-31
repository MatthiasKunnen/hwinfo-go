package hwinfoshmem

import (
	"github.com/MatthiasKunnen/hwinfo-go/pkg/util/bytesutil"
)

const (
	hwinfoSensorStringLength = 128
	hwinfoUnitStringLength   = 16
)

// HwinfoSensorStringAscii is a fixed length byte array of 8-bit ASCII encoded characters.
// The specific extended ASCII codepage used depends on the system's locale.
//
// Get the codepage used by your system using this powershell command:
//
//	[System.Text.Encoding]::Default
//
// The string it contains is padded by nul bytes.
type HwinfoSensorStringAscii [hwinfoSensorStringLength]byte

// HwinfoSensorStringUtf8 is a fixed length byte array of UTF-8 encoded characters.
// The string it contains is padded by nul bytes.
// To convert it to a string, use [HwinfoSensorStringUtf8.String].
// It is used in labels for sensor and reading.
type HwinfoSensorStringUtf8 [hwinfoSensorStringLength]byte

func (s HwinfoSensorStringUtf8) String() string {
	return bytesutil.Utf8BytesToString(s[:])
}

// HwinfoUnitStringAscii is the same as [HwinfoSensorStringAscii] but used for unit strings such as
// °C and MHz.
type HwinfoUnitStringAscii [hwinfoUnitStringLength]byte

// HwinfoUnitStringUtf8 is the same as  [HwinfoSensorStringUtf8] but used for unit strings such as
// °C and MHz.
type HwinfoUnitStringUtf8 [hwinfoUnitStringLength]byte

func (s HwinfoUnitStringUtf8) String() string {
	return bytesutil.Utf8BytesToString(s[:])
}
