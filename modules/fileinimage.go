package modules

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/klauspost/compress/zstd"
)

func Fileinimage_main(inFile string, outFile string, decode bool, noCompress bool) {
	if decode {
		fileinimage_decode(inFile, outFile, noCompress)
		fmt.Println("File decoded!")
		os.Exit(0)
	} else {
		fileinimage_encode(inFile, outFile, noCompress)
		fmt.Println("File encoded!")
		os.Exit(0)
	}
}

func fileinimage_encode(inFile string, outFile string, noCompress bool) {
	// Read the file
	data, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
	}

	var processedData []byte

	// If compression is enabled
	if !noCompress {
		// Compress the data with zstd
		w, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
		if err != nil {
			log.Fatal(err)
		}
		processedData = w.EncodeAll(data, nil)
	} else {
		processedData = data
	}

	// Calculate the width and height of the image
	size := int(math.Ceil(math.Sqrt(float64(len(processedData)))))

	// Create a new image
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Map each byte to a color and set it in the image
	for i, b := range processedData {
		x := i % size
		y := i / size
		img.Set(x, y, color.RGBA{R: b, G: b, B: b, A: 255})
	}

	// Add marker at the end
	img.Set(len(processedData)%size, len(processedData)/size, color.RGBA{R: 255, G: 10, B: 255, A: 100})

	// Create the output file
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Encode the image as a PNG
	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}

func fileinimage_decode(inFile string, outFile string, noCompress bool) {
	// Open the input file
	in, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	// Decode the image
	img, _, err := image.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	// Get the image bounds
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var data []byte

	// Extract the data from the image
dataLoop:
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r == 14109 && g == 2519 && b == 14109 && a == 25700 {
				fmt.Println("Found marker at", x, y, "Stopping decoding")
				// Stop decoding when we encounter the marker
				break dataLoop
			}
			data = append(data, byte(r))
		}
	}

	var processedData []byte

	// If compression was used
	if !noCompress {
		// Decompress the data with zstd
		d, err := zstd.NewReader(nil)
		if err != nil {
			log.Fatal(err)
		}
		processedData, err = d.DecodeAll(data, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		processedData = data
	}

	// Write the data to the output file
	err = os.WriteFile(outFile, processedData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
