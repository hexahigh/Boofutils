package modules

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"strings"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
	"github.com/mewkiz/flac"
)

//go:embed embed/audio/*
var audioFS embed.FS

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

func PlayAudioAdvanced(ctx context.Context, audioFiles string, bits int, channels int, sampleRate int) {
	// Split the audioFiles string into a slice of file names
	files := strings.Split(audioFiles, ",")

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.
	op := &oto.NewContextOptions{}
	op.SampleRate = sampleRate
	op.ChannelCount = channels
	switch bits {
	case 8:
		op.Format = oto.FormatUnsignedInt8
	case 16:
		op.Format = oto.FormatSignedInt16LE
	case 32:
		op.Format = oto.FormatFloat32LE
	}

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

			switch strings.Split(audioFile, ".")[1] {
			case "wav":
				// Skip the first 44 bytes as they are the wav header
				fileBytesReader.Seek(44, io.SeekStart)
				player = otoCtx.NewPlayer(fileBytesReader)
			case "mp3":
				decodedMp3, err := mp3.NewDecoder(fileBytesReader)
				if err != nil {
					panic("mp3.NewDecoder failed: " + err.Error())
				}
				player = otoCtx.NewPlayer(decodedMp3)
			case "flac":
				wavReader, err := flac2wav("embed/audio/" + audioFile)
				if err != nil {
					panic("flac2wav failed: " + err.Error())
				}
				player = otoCtx.NewPlayer(wavReader)
			default:
				panic("unsupported audio format: " + audioFile)
			}

			// Play starts playing the sound and returns without waiting for it (Play() is async).
			player.Play()

			// We can wait for the sound to finish playing using something like this
			for player.IsPlaying() {
				time.Sleep(time.Millisecond)
			}

			_, err = player.Seek(0, io.SeekStart)
			if err != nil {
				panic("player.Seek failed: " + err.Error())
			}
		}
	}
}

func flac2wav(path string) (io.Reader, error) {
	// Read the FLAC file from the embedded filesystem.
	flacData, err := fs.ReadFile(audioFS, path)
	if err != nil {
		return nil, err
	}

	// Create a buffer to hold the WAV data.
	var buf seekableBuffer

	// Create a FLAC stream from the embedded FLAC data.
	stream, err := flac.New(bytes.NewReader(flacData))
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	// Create WAV encoder.
	wavAudioFormat := 1 // PCM
	enc := wav.NewEncoder(&buf, int(stream.Info.SampleRate), int(stream.Info.BitsPerSample), int(stream.Info.NChannels), wavAudioFormat)
	defer enc.Close()
	var data []int
	for {
		// Decode FLAC audio samples.
		frame, err := stream.ParseNext()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Encode WAV audio samples.
		data = data[:0]
		for i := 0; i < frame.Subframes[0].NSamples; i++ {
			for _, subframe := range frame.Subframes {
				sample := int(subframe.Samples[i])
				if frame.BitsPerSample == 8 {
					const midpointValue = 0x80
					sample += midpointValue
				}
				data = append(data, sample)
			}
		}
		buf := &audio.IntBuffer{
			Format: &audio.Format{
				NumChannels: int(stream.Info.NChannels),
				SampleRate:  int(stream.Info.SampleRate),
			},
			Data:           data,
			SourceBitDepth: int(stream.Info.BitsPerSample),
		}
		if err := enc.Write(buf); err != nil {
			return nil, err
		}
	}

	return &buf, nil
}

// Custom type that embeds bytes.Buffer and implements io.Seeker.
type seekableBuffer struct {
	bytes.Buffer
}

// Implement the Seek method for the custom type.
func (sb *seekableBuffer) Seek(offset int64, whence int) (int64, error) {
	// For this use case, we only need to support seeking to the beginning.
	if whence == io.SeekStart && offset == 0 {
		sb.Buffer.Reset()
		return 0, nil
	}
	return 0, fmt.Errorf("only support seeking to the beginning")
}
