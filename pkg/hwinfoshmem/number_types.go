package hwinfoshmem

import (
	"encoding/binary"
	"math"
)

// HwinfoFloat64 represents Go's float64 and C#'s double.
// The type exists to deal with differences in memory alignment between Go and HWiNFO's shared
// memory.
type HwinfoFloat64 [8]byte

func (b HwinfoFloat64) ToFloat64() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(b[:]))
}
