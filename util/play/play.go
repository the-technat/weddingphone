package play

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
	"github.com/the-technat/weddingphone/util/errors"
)

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}

type ID [4]byte

// PlayMP3 will play the provided mp3 file
func PlayMP3(ctx context.Context, file string) error {
	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	errors.CheckError(err)

	errors.CheckError((decoder.Open(file)))
	defer decoder.Close()

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	out := make([]int16, 8192)
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	errors.CheckError(err)
	defer stream.Close()

	errors.CheckError(stream.Start())
	defer stream.Stop()
	for {
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err == mpg123.EOF {
			break
		}
		errors.CheckError(err)

		errors.CheckError((binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out)))
		errors.CheckError((stream.Write()))
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
	return nil
}

func PlayAIFF(ctx context.Context, file string) error {
	f, err := os.Open(file)
	errors.CheckError(err)
	defer f.Close()

	id, data, err := readChunk(f)
	errors.CheckError(err)
	if id.String() != "FORM" {
		return fmt.Errorf("bad file format")
	}
	_, err = data.Read(id[:])
	errors.CheckError(err)
	if id.String() != "AIFF" {
		return fmt.Errorf("bad file format")
	}
	var c commonChunk
	var audio io.Reader
	for {
		id, chunk, err := readChunk(data)
		if err == io.EOF {
			break
		}
		errors.CheckError(err)
		switch id.String() {
		case "COMM":
			errors.CheckError(binary.Read(chunk, binary.BigEndian, &c))
		case "SSND":
			chunk.Seek(8, 1) //ignore offset and block
			audio = chunk
		default:
			fmt.Printf("ignoring unknown chunk '%s'\n", id)
		}
	}

	// assume 44100 sample rate, mono, 32 bit

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)
	errors.CheckError(err)
	defer stream.Close()

	errors.CheckError(stream.Start())
	defer stream.Stop()
	for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		}
		errors.CheckError(err)
		errors.CheckError(stream.Write())
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
	return nil
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

func (id ID) String() string {
	return string(id[:])
}
