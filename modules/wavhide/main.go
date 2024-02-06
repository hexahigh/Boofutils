package wavhide

import (
	"fmt"
	"log"
	"os"
)

func Main(inputFile string, wavFile string, outputFile string) {

	// Read the input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	// Process each byte
	processedData := make([]byte, len(data))
	for i, b := range data {
		processedData[i] = b ^ byte(10000&0xFF) // XOR with the lower  8 bits of  10000
	}

	// Read the WAV file
	wavData, err := os.ReadFile(wavFile)
	if err != nil {
		log.Fatalf("Error reading WAV file: %v\n", err)
	}

	// Append the processed data to the WAV file
	appendedData := append(wavData, processedData...)

	// Write the appended data to the output file
	err = os.WriteFile(outputFile, appendedData, 0644)
	if err != nil {
		log.Fatalf("Error writing output WAV file: %v\n", err)
	}

	fmt.Println("Process completed successfully.")
}
