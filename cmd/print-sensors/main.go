package main

import (
	"github.com/MatthiasKunnen/hwinfo-go/internal/text"
	"github.com/MatthiasKunnen/hwinfo-go/pkg/hwinfoshmem"
)
import (
	"fmt"
	"os"
)

func main() {
	var memoryReader = hwinfoshmem.NewMemoryReader()

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
			reading.UserLabelUtf8.String(),
			fmt.Sprintf("%f", reading.GetValue()),
			reading.UnitUtf8.String(),
		})
	}

	err = printer.Write()
	if err != nil {
		fmt.Printf("Error printing table: %s\n", err)
		os.Exit(1)
	}
	copyReader := memoryReader.Copy(hwInfo)

	os.WriteFile("memcopy.bin", copyReader.Bytes, 0666)
}
