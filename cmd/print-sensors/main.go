package main

import "github.com/MatthiasKunnen/hwinfo-go/pkg/sharedmemory"
import (
	"fmt"
	"os"
)

func main() {
	var memoryReader = sharedmemory.NewMemoryReader()

	err := memoryReader.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = memoryReader.Lock()
	defer memoryReader.ReleaseLock()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hwInfo, err := memoryReader.Reader.GetHeader()
	if err != nil {
		fmt.Printf("Failed to get header: %s\n", err)
		os.Exit(1)
	}

	if !hwInfo.IsActive() {
		fmt.Println("HWiNFO is not active")
		os.Exit(1)
	}

	readings, err := memoryReader.Reader.GetReadings(hwInfo)
	if err != nil {
		fmt.Printf("Error getting readings %v\n", err)
		os.Exit(1)
	}

	for _, reading := range readings {
		fmt.Printf("%s\t(%d/%d)\t%f\t%s\n", reading.GetUserLabel(), reading.SensorIndex, reading.Id, reading.GetValue(), reading.GetUnit())
	}

	copyReader := memoryReader.Copy(hwInfo)
	copiedReadings, err := copyReader.Reader.GetReadings(hwInfo)
	if err != nil {
		fmt.Printf("Error getting copied readings %v\n", err)
		os.Exit(1)
	}

	for _, reading := range copiedReadings {
		fmt.Printf("%s\t(%d/%d)\t%f\t%s\n", reading.GetUserLabel(), reading.SensorIndex, reading.Id, reading.GetValue(), reading.GetUnit())
	}

	readingsById, err := memoryReader.Reader.GetReadingsById(hwInfo, []sharedmemory.ReadingIdSensorCombo{
		{
			Id:          134217730,
			SensorIndex: 24,
		},
	})
	if err != nil {
		fmt.Printf("Error getting readings by ID %v\n", err)
		os.Exit(1)
	}

	for _, reading := range readingsById {
		fmt.Printf("By ID: %s\t(%d/%d)\t%f\t%s\n", reading.GetUserLabel(), reading.SensorIndex, reading.Id, reading.GetValue(), reading.GetUnit())
	}
}
