package modules

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
)

func pixelToAnsi(px1, px2 color.Color) string {
	res := ""
	r1, g1, b1, a1 := px1.RGBA()
	r2, g2, b2, a2 := px2.RGBA()

	if a1 != 65535 || a2 != 65535 {
		res += "\033[0m"
	}
	if a1 != 65535 && a2 != 65535 {
		res += " "
	}
	if a1 == 65535 {
		res += fmt.Sprintf("\033[38;2;%d;%d;%dm", r1/256, g1/256, b1/256)
		if a2 == 65535 {
			res += fmt.Sprintf("\033[48;2;%d;%d;%dm", r2/256, g2/256, b2/256)
		}
		res += "▀"
	} else if a2 == 65535 {
		res += fmt.Sprintf("\033[38;2;%d;%d;%dm", r2/256, g2/256, b2/256)
		res += "▄"
	}
	return res
}

func imageToAnsi(img image.Image) string {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	finalRes := ""

	for i := 0; i < height; i += 2 {
		for j := 0; j < width; j++ {
			finalRes += pixelToAnsi(img.At(j, i), img.At(j, i+1))
		}
		finalRes += "\033[0m\n"
	}
	if height%2 != 0 {
		for j := 0; j < width; j++ {
			finalRes += pixelToAnsi(img.At(j, height-1), color.Transparent)
		}
	}

	return finalRes
}

func Ansiimg_main(filename string, output string, width uint, height uint) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	// Resize the image
	if height == 0 {
		img = resize.Resize(width, 0, img, resize.Lanczos3)
	} else {
		img = resize.Resize(width, height, img, resize.Lanczos3)
	}

	finalRes := imageToAnsi(img)

	if output != "" {
		outFile, err := os.Create(output)
		if err != nil {
			fmt.Println("Error: File could not be created")
			os.Exit(1)
		}
		defer outFile.Close()

		outFile.WriteString(finalRes)
	} else {
		fmt.Println(finalRes)
	}
}
