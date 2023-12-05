package modules

import (
	"fmt"
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/klauspost/compress/zstd"
)

func Fileinaudio_main(inFile string, outFile string, decode bool, noCompress bool) {
	if decode {
		fileinaudio_decode(inFile, outFile, noCompress)
		fmt.Println("File decoded!")
		os.Exit(0)
	} else {
		fileinaudio_encode(inFile, outFile, noCompress)
		fmt.Println("File encoded!")
		os.Exit(0)
	}
}

func fileinaudio_encode(inFile string, outFile string, noCompress bool) {
	var buf *audio.IntBuffer

	// Read the file
	data, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
	}

	// If compression is enabled
	if !noCompress {
		// Compress the data with zstd
		w, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(22)))
		if err != nil {
			log.Fatal(err)
		}
		compressedData := w.EncodeAll(data, nil)

		// Create a new audio.IntBuffer
		buf = &audio.IntBuffer{Data: make([]int, len(compressedData)), Format: &audio.Format{SampleRate: 44100, NumChannels: 1}}

		// Map each byte to a frequency and add it to the buffer
		for i, b := range compressedData {
			buf.Data[i] = int(b) * 100 // Multiply by 100 to get a frequency in the audible range
		}
		// If compression is not enabled
	} else {
		// Create a new audio.IntBuffer
		buf = &audio.IntBuffer{Data: make([]int, len(data)), Format: &audio.Format{SampleRate: 44100, NumChannels: 1}}
		// Map each byte to a frequency and add it to the buffer
		for i, b := range data {
			buf.Data[i] = int(b) * 100 // Multiply by 100 to get a frequency in the audible range
		}
	}

	// Create a new wav.Encoder
	out, err := os.Create(outFile)
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
}

func fileinaudio_decode(inFile string, outFile string, noCompress bool) {
	// Open the .wav file
	in, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	dec := wav.NewDecoder(in)

	// Read the .wav file
	buf, err := dec.FullPCMBuffer()
	if err != nil {
		log.Fatal(err)
	}

	// Create a byte slice to hold the decoded data
	data := make([]byte, len(buf.Data))

	// Map each frequency to a byte and add it to the slice
	for i, f := range buf.Data {
		data[i] = byte(f / 100) // Divide by 100 to get the original byte
	}

	if !noCompress {
		// Decompress the data with zstd
		d, err := zstd.NewReader(nil)
		if err != nil {
			log.Fatal(err)
		}
		data, err = d.DecodeAll(data, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Write the data to a file
	err = os.WriteFile(outFile, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
