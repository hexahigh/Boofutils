package modules

import (
	"archive/tar"
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	gzip "github.com/klauspost/pgzip"
)

//go:embed embed/audio/*
var audioFS embed.FS

func Bua_main(inFile string, outFile string, encode bool, mute bool, v int) {
	if encode {
		Bua_encode(inFile, outFile, mute, v)
	} else {
		Bua_decode(inFile, outFile, mute, v)
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

func Bua_encode(inFile string, outFile string, mute bool, v int) {
	ctx, cancel := context.WithCancel(context.Background())
	if !mute {
		go PlayAudioMult(ctx, "carry_on.mp3")
	}
	files := strings.Split(inFile, ",")
	tarfile, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer tarfile.Close()
	bw := gzip.NewWriter(tarfile)
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

// ! DEPRECATED! Use PlayAudioMult instead
func PlayAudioLoop(ctx context.Context, audioFile string) {
	// Read the mp3 file into memory
	fileBytes, err := audioFS.ReadFile("embed/audio/" + audioFile)
	if err != nil {
		panic("reading " + audioFile + " failed: " + err.Error())
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

	for {
		select {
		case <-ctx.Done():
			return
		default:
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
		}
	}
}

/*
Decodes and plays a random audio file from a comma separated list.
Audio files must be located in the modules/embed/audio directory.
Supports wav and mp3.
*/
func PlayAudioMult(ctx context.Context, audioFiles string) {
	// Split the audioFiles string into a slice of file names
	files := strings.Split(audioFiles, ",")

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.
	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	<-readyChan

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Choose a random file from the slice
			audioFile := files[rand.Intn(len(files))]

			// Read the audio file into memory
			fileBytes, err := audioFS.ReadFile("embed/audio/" + audioFile)
			if err != nil {
				panic("reading " + audioFile + " failed: " + err.Error())
			}

			// Convert the bytes into a reader object that can be used with the decoder
			fileBytesReader := bytes.NewReader(fileBytes)

			var player *oto.Player

			// Check the file extension and use the appropriate decoder
			if strings.HasSuffix(audioFile, ".wav") {
				// Skip the first 44 bytes as they are the wav header
				fileBytesReader.Seek(44, io.SeekStart)
				player = otoCtx.NewPlayer(fileBytesReader)
			} else if strings.HasSuffix(audioFile, ".mp3") {
				decodedMp3, err := mp3.NewDecoder(fileBytesReader)
				if err != nil {
					panic("mp3.NewDecoder failed: " + err.Error())
				}
				player = otoCtx.NewPlayer(decodedMp3)
			} else {
				panic("unsupported audio format: " + audioFile)
			}

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
		}
	}
}
