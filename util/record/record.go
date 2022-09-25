package record

import (
	"encoding/binary"
	"os"

	"github.com/gordonklaus/portaudio"
	"github.com/the-technat/weddingphone/util/errors"
)

func RecordToFile(file string, stop chan os.Signal) error {
	f, err := os.Create(file)
	errors.CheckError(err)

	// form chunk
	_, err = f.WriteString("FORM")
	errors.CheckError(err)
	errors.CheckError(binary.Write(f, binary.BigEndian, int32(0))) //total bytes
	_, err = f.WriteString("AIFF")
	errors.CheckError(err)

	// common chunk
	_, err = f.WriteString("COMM")
	errors.CheckError(err)
	errors.CheckError(binary.Write(f, binary.BigEndian, int32(18)))    //size
	errors.CheckError(binary.Write(f, binary.BigEndian, int16(1)))     //channels
	errors.CheckError(binary.Write(f, binary.BigEndian, int32(0)))     //number of samples
	errors.CheckError(binary.Write(f, binary.BigEndian, int16(32)))    //bits per sample
	_, err = f.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0}) //80-bit sample rate 44100
	errors.CheckError(err)

	// sound chunk
	_, err = f.WriteString("SSND")
	errors.CheckError(err)
	errors.CheckError(binary.Write(f, binary.BigEndian, int32(0))) //size
	errors.CheckError(binary.Write(f, binary.BigEndian, int32(0))) //offset
	errors.CheckError(binary.Write(f, binary.BigEndian, int32(0))) //block
	nSamples := 0
	defer func() {
		// fill in missing sizes
		totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples
		_, err = f.Seek(4, 0)
		errors.CheckError(err)
		errors.CheckError(binary.Write(f, binary.BigEndian, int32(totalBytes)))
		_, err = f.Seek(22, 0)
		errors.CheckError(err)
		errors.CheckError(binary.Write(f, binary.BigEndian, int32(nSamples)))
		_, err = f.Seek(42, 0)
		errors.CheckError(err)
		errors.CheckError(binary.Write(f, binary.BigEndian, int32(4*nSamples+8)))
		errors.CheckError(f.Close())
	}()

	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(2, 0, 44100, len(in), in)
	errors.CheckError(err)
	defer stream.Close()

	errors.CheckError(stream.Start())
	for {
		errors.CheckError(stream.Read())
		errors.CheckError(binary.Write(f, binary.BigEndian, in))
		nSamples += len(in)
		select {
		case <-stop:
			errors.CheckError(stream.Stop())
		default:
		}
	}
}
