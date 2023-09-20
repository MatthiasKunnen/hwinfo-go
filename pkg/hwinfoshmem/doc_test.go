package hwinfoshmem_test

import (
	_ "embed"
	"fmt"
	"github.com/MatthiasKunnen/hwinfo-go/internal/text"
	"github.com/MatthiasKunnen/hwinfo-go/pkg/hwinfoshmem"
	"os"
)

//go:embed testdata/limited_live.bin
var data []byte

// Example of the structs and their functions provided by the library.
func Example() {
	// For demo purposes, uses copy of shared memory. See MemoryReader for reading the actual shared memory.
	var memoryReader = hwinfoshmem.NewBytesReader(data)

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

	printer := text.NewTablePrinter(os.Stdout, make([]text.Column, 3), "    ")

	printer.Append([]string{"Label", "Value", "Unit"})

	for _, reading := range readings {
		printer.Append([]string{
			reading.UserLabel.String(),
			fmt.Sprintf("%f", reading.Value.ToFloat64()),
			reading.Unit.String(),
		})
	}

	fmt.Println("Limited sensors:")
	err = printer.Write()
	if err != nil {
		fmt.Printf("Failed to print table: %s\n", err)
	}

	// Output: Limited sensors:
	// Label                              Value        Unit
	// CPU (Tctl/Tdie)                    47.250000    °C
	// CPU Die (average)                  45.087887    °C
	// CPU CCD1 (Tdie)                    45.125000    °C
	// CPU CCD2 (Tdie)                    33.375000    °C
	// Water (EC_TEMP1)                   27.000000    °C
	// GPU Memory Junction Temperature    48.000000    °C
	// GPU Hot Spot Temperature           35.000000    °C
}
