package main

import "github.com/MatthiasKunnen/hwinfo-go/pkg/sharedmemory"
import (
	"fmt"
	"os"
)

func main() {
	var reader = sharedmemory.NewMemoryReader()

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

	hwInfo, err := reader.Reader.GetHeader()
	if err != nil {
		fmt.Printf("Failed to get header: %s\n", err)
		os.Exit(1)
	}

	if !hwInfo.IsActive() {
		fmt.Println("HWiNFO is not active")
		os.Exit(1)
	}

	readings, err := reader.Reader.GetReadings(hwInfo)
	if err != nil {
		fmt.Printf("Error getting readings %v\n", err)
		os.Exit(1)
	}

	for _, reading := range readings {
		fmt.Printf("%s\t(%d/%d)\t%f\t%s\n", reading.GetUserLabel(), reading.SensorIndex, reading.Id, reading.GetValue(), reading.GetUnit())
	}

	copyReader := reader.Copy(hwInfo)
	copiedReadings, err := copyReader.Reader.GetReadings(hwInfo)
	if err != nil {
		fmt.Printf("Error getting copied readings %v\n", err)
		os.Exit(1)
	}

	for _, reading := range copiedReadings {
		fmt.Printf("%s\t(%d/%d)\t%f\t%s\n", reading.GetUserLabel(), reading.SensorIndex, reading.Id, reading.GetValue(), reading.GetUnit())
	}

	readingsById, err := reader.Reader.GetReadingsById(hwInfo, []sharedmemory.ReadingIdSensorCombo{
		{
			Id:          134217730,
			SensorIndex: 24,
		},
	})
	if err != nil {
		fmt.Printf("Error getting copied readings %v\n", err)
		os.Exit(1)
	}

	for _, reading := range readingsById {
		fmt.Printf("By ID: %s\t(%d/%d)\t%f\t%s\n", reading.GetUserLabel(), reading.SensorIndex, reading.Id, reading.GetValue(), reading.GetUnit())
	}
}
