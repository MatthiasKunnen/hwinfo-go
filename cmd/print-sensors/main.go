package main

import (
	"github.com/MatthiasKunnen/hwinfo-go/internal/text"
	"github.com/MatthiasKunnen/hwinfo-go/pkg/sharedmemory"
)
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

	hwInfo, err := memoryReader.Data.GetHeader()
	if err != nil {
		fmt.Printf("Failed to get header: %s\n", err)
		os.Exit(1)
	}

	if !hwInfo.IsActive() {
		fmt.Println("HWiNFO is not active")
		os.Exit(1)
	}

	readings, err := memoryReader.Data.GetReadings(hwInfo)
	if err != nil {
		fmt.Printf("Error getting readings %v\n", err)
		os.Exit(1)
	}

	printer := text.NewTablePrinter(os.Stdout, make([]text.Column, 3), "    ")

	printer.Append([]string{"Label", "Value", "Unit"})

	for _, reading := range readings {
		printer.Append([]string{
			reading.GetUserLabel(),
			fmt.Sprintf("%f", reading.GetValue()),
			reading.GetUnit(),
		})
	}

	err = printer.Write()
	if err != nil {
		fmt.Printf("Error printing table: %s\n", err)
		os.Exit(1)
	}
	copyReader := memoryReader.Copy(hwInfo)
	copiedReadings, err := copyReader.Data.GetReadings(hwInfo)
	if err != nil {
		fmt.Printf("Error getting copied readings %v\n", err)
		os.Exit(1)
	}

	for _, reading := range copiedReadings {
		fmt.Printf("%s\t(%d/%d)\t%f\t%s\n", reading.GetUserLabel(), reading.SensorIndex, reading.Id, reading.GetValue(), reading.GetUnit())
	}

	readingsById, err := memoryReader.Data.GetReadingsById(hwInfo, []sharedmemory.ReadingIdSensorCombo{
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
