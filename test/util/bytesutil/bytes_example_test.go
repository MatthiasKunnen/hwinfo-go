package bytesutil_test

import (
	"fmt"
	"github.com/MatthiasKunnen/hwinfo-go/internal/text"
	"github.com/MatthiasKunnen/hwinfo-go/pkg/util/bytesutil"
	"os"
	"unicode/utf8"
)

func ExampleUtf8BytesToString() {
	theBytes := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0, 0, 0}

	printer := text.NewTablePrinter(os.Stdout, make([]text.Column, 3), "    ")

	printer.Append([]string{"function", "Result (%q)", "Amount of characters"})

	normalToString := string(theBytes[:])
	printer.Append([]string{
		"string(byteArray[:])",
		fmt.Sprintf("%q", normalToString),
		fmt.Sprintf("%d", utf8.RuneCountInString(normalToString)),
	})

	nulTerminatedResult := bytesutil.Utf8BytesToString(theBytes)
	printer.Append([]string{
		"Utf8BytesToString(byteArray)",
		fmt.Sprintf("%q", nulTerminatedResult),
		fmt.Sprintf("%d", utf8.RuneCountInString(nulTerminatedResult)),
	})

	fmt.Printf("byteArray = % 02x\n", theBytes)
	fmt.Printf("\nComparison of byte to string conversions:\n")
	err := printer.Write()
	if err != nil {
		fmt.Printf("error printing table: %s\n", err)
	}

	// Output:
	// byteArray = 48 65 6c 6c 6f 00 00 00
	//
	// Comparison of byte to string conversions:
	// function                        Result (%q)            Amount of characters
	// string(byteArray[:])            "Hello\x00\x00\x00"    8
	// Utf8BytesToString(byteArray)    "Hello"                5
}
