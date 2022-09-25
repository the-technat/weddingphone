package play

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
	"github.com/the-technat/weddingphone/util/errors"
)

// PlayMP3 will play the provided mp3 file
func PlayMP3(file string, close chan os.Signal) error {
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
		case <-close:
			return nil
		default:
		}
	}
	return nil
}
