package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func hex_main() {
	fmt.Println("What is the path of the file?")
	filePath := askInput()
	readLinesHexraw(filePath)
}

func readLinesHexraw(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// Read the file in chunks
	buf := make([]byte, 256) // Adjust the chunk size as needed
	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		// Convert the chunk to hexadecimal and print it
		fmt.Printf("%s", hex.EncodeToString(buf))
	}
	return nil, nil
}
