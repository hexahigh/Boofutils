package modules

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

func Bua_main(inFile string, outFile string, encode bool) {
	if encode {
		Bua_encode(inFile, outFile)
	} else {
		Bua_decode(inFile, outFile)
	}
}

func Bua_decode(inFile string, outDir string) {
	// Open the zstd compressed file
	zr, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()

	// Create a zstd decoder
	dec, err := zstd.NewReader(zr)
	if err != nil {
		log.Fatal(err)
	}
	defer dec.Close()

	// Create a tar reader
	tr := tar.NewReader(dec)

	// Iterate over the files in the tar archive
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}

		// The target location where to decompress the file
		target := filepath.Join(outDir, header.Name)

		// Check the file type
		switch header.Typeflag {
		case tar.TypeDir: // if a dir
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				log.Fatal(err)
			}
		case tar.TypeReg: // if a file
			// Ensure the parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				log.Fatal(err)
			}

			// Create the file
			f, err := os.Create(target)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			// Copy data from the tar archive to the file
			if _, err := io.Copy(f, tr); err != nil {
				log.Fatal(err)
			}
		default:
			log.Printf("Can't: %c, %s\n", header.Typeflag, target)
		}
	}
}

func Bua_encode(inFile string, outFile string) {
	// Split the inFile string into a slice of file paths
	files := strings.Split(inFile, ",")

	tarfile, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer tarfile.Close()

	// create a new zstd writer
	zw, err := zstd.NewWriter(tarfile, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		log.Fatal(err)
	}
	defer zw.Close()

	tw := tar.NewWriter(zw)
	defer tw.Close()

	// Iterate over the files and add them to the tar archive
	for _, file := range files {
		file = strings.TrimSpace(file) // Remove any leading/trailing white space
		err = filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(info, path)
			if err != nil {
				return err
			}

			header.Name = path // Ensure the name is correct
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			if !info.Mode().IsRegular() { // Skip if not a regular file
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(tw, f)
			return err
		})

		if err != nil {
			log.Fatal(err)
		}
	}
}
