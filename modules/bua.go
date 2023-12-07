package modules

import (
	"archive/tar"
	"bytes"
	"embed"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"github.com/klauspost/compress/zstd"
)

//go:embed embed/audio_test.mp3
var audioF embed.FS

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

func PlayAudio() {
	// Read the mp3 file into memory
	fileBytes, err := audioF.ReadFile("embed/audio_test.mp3")
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	// Create a new 'player' that will handle our sound. Paused by default.
	player := otoCtx.NewPlayer(decodedMp3)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	newPos, err := player.Seek(0, io.SeekStart)
	if err != nil {
		panic("player.Seek failed: " + err.Error())
	}
	println("Player is now at position:", newPos)
	player.Play()

}
