package hwinfoshmem_test

import (
	"fmt"
	"github.com/MatthiasKunnen/hwinfo-go/pkg/hwinfoshmem"
)

func ExampleMemoryReader() {
	var memoryReader = hwinfoshmem.NewMemoryReader()

	err := memoryReader.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = memoryReader.Lock()
	defer memoryReader.ReleaseLock()
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

	for _, reading := range readings {
		fmt.Printf("%-50s\t(%d/%d)\t%f\t%s\n", reading.UserLabel, reading.SensorIndex, reading.Id, reading.Value.ToFloat64(), reading.Unit)
	}
}
