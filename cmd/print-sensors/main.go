package main

import "github.com/MatthiasKunnen/hwinfo-go/pkg/sharedmemory"
import (
	"fmt"
	"os"
)

func main() {
	var reader = sharedmemory.MemoryReader{}

	err := reader.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = reader.Lock()
	defer reader.ReleaseLock()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hwInfo, err := reader.GetHeader()
	if err != nil {
		fmt.Printf("Failed to get header: %s\n", err)
		os.Exit(1)
	}

	if !hwInfo.IsActive() {
		fmt.Println("HWiNFO is not active")
		os.Exit(1)
	}

	readings, err := reader.GetReadings(hwInfo)
	if err != nil {
		fmt.Printf("Error getting readings %v\n", err)
		os.Exit(1)
	}

	for _, reading := range readings {
		fmt.Printf("%s\t(%d/%d)\t%f\t%s\n", reading.GetUserLabel(), reading.SensorIndex, reading.Id, reading.GetValue(), reading.GetUnit())
	}
}
