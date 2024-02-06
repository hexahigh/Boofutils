package wavhide

import (
	"fmt"
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func Main(inputFile string, wavFile string, outputFile string) {
	var buf *audio.IntBuffer
	// Open the input file
	inputFd, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error opening input file: %v\n", err)
	}
	// Create a new audio.IntBuffer
	buf = &audio.IntBuffer{Data: make([]int, len(inputFd)), Format: &audio.Format{SampleRate: 44100, NumChannels: 1}}
	// Map each byte to a frequency and add it to the buffer
	for i, b := range inputFd {
		buf.Data[i] = int(b) * 100 // Multiply by 100 to get a frequency in the audible range
	}

	// Create a new wav.Encoder
	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	enc := wav.NewEncoder(out, buf.Format.SampleRate, 16, buf.Format.NumChannels, 1)

	// Write the buffer to the encoder
	if err := enc.Write(buf); err != nil {
		log.Fatal(err)
	}

	// Close the encoder
	if err := enc.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Process completed successfully.")
}
