package hwinfoshmem

import (
	"encoding/binary"
	"github.com/MatthiasKunnen/hwinfo-go/pkg/util"
	"time"
)

// Do not use any other type than uint32, int32, and [?]byte to avoid memory alignment issues

// HwinfoHeader contains information regarding the rest of the available shared memory.
type HwinfoHeader struct {
	// Reports whether HWiNFO is active.
	// "HWiS" when HWiNFO it is Active, "DAED" (sic.) when it is not.
	// Get string using GetStatus or check if it is live using isActive.
	Status [4]byte

	// Structure layout version. 1=Initial; 2=Added UTF-8 strings (HWiNFO v7.33+)
	Version uint32

	// Options:
	//  - 0: Initial layout (HWiNFO ver <= 6.11)
	//  - 1: Added (HWiNFO v6.11-3917)
	Revision uint32

	// The unix time (seconds since 1970-01-01) when the last update to the data occurred.
	// Get int using GetLastUpdate.
	LastUpdate [8]byte

	// Offset of the Sensor section from beginning of HwinfoHeader.
	SensorSectionOffset uint32

	// The size in bytes of every sensor's data. sizeof(HwinfoSensor).
	SensorSize uint32

	// Amount of sensors that are available.
	SensorAmount uint32

	// Offset of the Reading section from beginning of HwinfoHeader.
	ReadingSectionOffset uint32

	// Size of each reading's data in bytes. sizeof(HwinfoReading)
	ReadingSize uint32

	// Number of readings.
	ReadingAmount uint32

	// Time in milliseconds between updates of the data by HWiNFO.
	PollingPeriodInMs uint32
}

// IsActive returns true when HWiNFO is currently updating the shared memory.
// When HWiNFO shared memory is not active, this usually means that the shared memory time limit
// has expired.
func (info HwinfoHeader) IsActive() bool {
	return info.Status == [4]byte{0x48, 0x57, 0x69, 0x53} // HWiS in ASCII
}

// GetStatus returns the status of the shared memory.
// "HWiS" when HWiNFO it is Active,
// "DAED" (sic.) when it is not.
// When HWiNFO shared memory is not active, this usually means that the shared memory time limit
// has expired.
func (info HwinfoHeader) GetStatus() string {
	return util.Utf8BytesToString(info.Status[:])
}

// GetLastUpdate returns the time since HWiNFO last updated the shared memory in seconds since
// Unix Epoch.
func (info HwinfoHeader) GetLastUpdate() int64 {
	return int64(binary.LittleEndian.Uint64(info.LastUpdate[:]))
}

// GetLastUpdateTime returns a Time object representing the time since HWiNFO updated the shared
// memory.
func (info HwinfoHeader) GetLastUpdateTime() time.Time {
	return time.Unix(info.GetLastUpdate(), 0)
}
