package modules

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	f "github.com/hexahigh/boofutils/modules/flagmanager"
	gzip "github.com/klauspost/pgzip"

	bzip2Decomp "github.com/cosnicolaou/pbzip2"
	bzip2Comp "github.com/dsnet/compress/bzip2"
)

func Bua_main(config f.BuaConfig) {
	l := gzip.DefaultCompression
	if !config.V2 {
		if config.BestCompression {
			l = gzip.BestCompression
		}
	} else {
		if config.BestCompression {
			l = bzip2Comp.BestCompression
		} else {
			l = bzip2Comp.DefaultCompression
		}
	}

	useV2 := config.V2

	// Add .bua to the output file if it has no extension
	if path.Ext(config.OutFile) == "" && config.Encode && !useV2 {
		config.OutFile = config.OutFile + ".bua"
	}
	if path.Ext(config.OutFile) == "" && config.Encode && useV2 {
		config.OutFile = config.OutFile + ".bua2"
	}
	if config.OutFile == "" && !config.Encode {
		config.OutFile = "./"
	}

	VerbPrintln(config.Verbosity, 1, "Starting with these options:")
	VerbPrintln(config.Verbosity, 1, config)

	if useV2 {
		if config.Encode {
			Bua_encode2(config.InFile, config.OutFile, config.Mute, config.Verbosity, l)
		} else {
			Bua_decode2(config.InFile, config.OutFile, config.Mute, config.Verbosity)
		}
	} else {
		if config.Encode {
			Bua_encode(config.InFile, config.OutFile, config.Mute, config.Verbosity, l)
		} else {
			Bua_decode(config.InFile, config.OutFile, config.Mute, config.Verbosity)
		}
	}
}

func Bua_decode(inFile string, outDir string, mute bool, v int) {
	// Start the music and console logging
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if !mute {
		go PlayAudioMult(ctx, "carry_on.mp3")
	}
	if outDir == "" {
		outDir = "."
	}
	if inFile == "" {
		VerbPrintln(v, 0, "No archive specified")
		VerbPrintln(v, 0, "Enter the path to the archive: ")
		inFile = AskInput()
	}

	// Open the bzip2 compressed file
	br, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer br.Close()

	// Create a bzip2 reader
	dec, err := gzip.NewReader(br)
	if err != nil {
		log.Fatal(err)
	}

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
			VerbPrintln(v, 0, "Extracted: ", target)
		default:
			log.Printf("Can't: %c, %s\n", header.Typeflag, target)
		}
	}
	VerbPrintln(v, 0, "Done!")
	VerbPrintln(v, 0, "Press any key to exit")
	fmt.Scanln()
	cancel()
}

func Bua_encode(inFile string, outFile string, mute bool, v int, l int) {
	ctx, cancel := context.WithCancel(context.Background())
	if !mute {
		go PlayAudioMult(ctx, "carry_on.mp3")
	}
	files := strings.Split(inFile, ",")
	tarfile, err := os.Create(outFile)
	if err != nil {
		VerbPrintln(v, 0, err.Error())
	}
	defer tarfile.Close()
	bw, err := gzip.NewWriterLevel(tarfile, l)
	if err != nil {
		VerbPrintln(v, 0, err.Error())
	}
	defer bw.Close()
	tw := tar.NewWriter(bw)
	defer tw.Close()
	// Iterate over the files and add them to the tar archive
	for _, file := range files {
		file = strings.TrimSpace(file) // Remove any leading/trailing white space
		baseDir := filepath.Dir(file)
		err = filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}
			VerbPrintln(v, 0, "Adding: ", relPath, "(", FileSize(path), ")")
			header, err := tar.FileInfoHeader(info, relPath)
			if err != nil {
				return err
			}
			header.Name = relPath // Ensure the name is correct
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
	cancel()
}

func Bua_decode2(inFile string, outDir string, mute bool, v int) {
	// Start the music and console logging
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if !mute {
		go PlayAudioMult(ctx, "carry_on.mp3")
	}
	if outDir == "" {
		outDir = "."
	}
	if inFile == "" {
		VerbPrintln(v, 0, "No archive specified")
		VerbPrintln(v, 0, "Enter the path to the archive: ")
		inFile = AskInput()
	}

	// Open the bzip2 compressed file
	br, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer br.Close()

	// Create a bzip2 reader
	dec := bzip2Decomp.NewReader(ctx, br)
	if err != nil {
		log.Fatal(err)
	}

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
			VerbPrintln(v, 0, "Extracted: ", target)
		default:
			log.Printf("Can't: %c, %s\n", header.Typeflag, target)
		}
	}
	VerbPrintln(v, 0, "Done!")
	VerbPrintln(v, 0, "Press any key to exit")
	fmt.Scanln()
	cancel()
}

func Bua_encode2(inFile string, outFile string, mute bool, v int, l int) {
	ctx, cancel := context.WithCancel(context.Background())
	if !mute {
		go PlayAudioMult(ctx, "carry_on.mp3")
	}
	files := strings.Split(inFile, ",")
	tarfile, err := os.Create(outFile)
	if err != nil {
		VerbPrintln(v, 0, err.Error())
	}
	defer tarfile.Close()
	bzconfig := &bzip2Comp.WriterConfig{Level: l}
	bw, err := bzip2Comp.NewWriter(tarfile, bzconfig)
	if err != nil {
		VerbPrintln(v, 0, err.Error())
	}
	defer bw.Close()
	tw := tar.NewWriter(bw)
	defer tw.Close()
	// Iterate over the files and add them to the tar archive
	for _, file := range files {
		file = strings.TrimSpace(file) // Remove any leading/trailing white space
		baseDir := filepath.Dir(file)
		err = filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}
			VerbPrintln(v, 0, "Adding: ", relPath, "(", FileSize(path), ")")
			header, err := tar.FileInfoHeader(info, relPath)
			if err != nil {
				return err
			}
			header.Name = relPath // Ensure the name is correct
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
	cancel()
}
