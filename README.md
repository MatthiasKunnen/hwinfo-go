# WIP: Go library for interfacing with [HWiNFO](https://www.hwinfo.com/)

Supports reading [HWiNFO](https://www.hwinfo.com/)'s Shared Memory.
Use cases:
- Make your own UI to display specific sensor values
- Execute some code if a sensor value is exceeded
- Log sensor values

## Documentation
- Shared memory: <https://pkg.go.dev/github.com/MatthiasKunnen/hwinfo-go/pkg/hwinfoshmem>

## Examples

### Print all HWiNFO readings

```go
package main

import (
	"fmt"
	"github.com/MatthiasKunnen/hwinfo-go/pkg/hwinfoshmem"
)

func main() {
	var memoryReader = hwinfoshmem.NewMemoryReader()

	err := memoryReader.Open()
	defer memoryReader.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = memoryReader.Lock()
	if err != nil {
		fmt.Println(err)
		return
	}

	hwInfo, err := memoryReader.GetHeader()
	if err != nil {
		fmt.Printf("Failed to get header: %s\n", err)
		return
	}

	if !hwInfo.IsActive() {
		fmt.Println("HWiNFO is not active")
		return
	}

	readings, err := memoryReader.GetReadings(hwInfo)
	if err != nil {
		fmt.Printf("Error getting readings %v\n", err)
		return
	}

	fmt.Printf("%-35s\t%s\t%s\n", "Label", "Value", "Unit")
	for _, reading := range readings {
		fmt.Printf("%-35s\t%f\t%s\n", reading.UserLabel, reading.Value.ToFloat64(), reading.Unit)
	}
}
```

Outputs
```
Label                              Value        Unit
CPU (Tctl/Tdie)                    47.250000    °C
CPU Die (average)                  45.087887    °C
CPU CCD1 (Tdie)                    45.125000    °C
CPU CCD2 (Tdie)                    33.375000    °C
Water (EC_TEMP1)                   27.000000    °C
GPU Memory Junction Temperature    48.000000    °C
GPU Hot Spot Temperature           35.000000    °C
...
```
