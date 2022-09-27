package main

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/gordonklaus/portaudio"
	log "github.com/sirupsen/logrus"
	"github.com/the-technat/weddingphone/util/play"
	"github.com/the-technat/weddingphone/util/record"
)

const (
	// INTRO - the file that should be played as intro to the users
	INTRO = "/home/technat/code/weddingphone/assets/intro.aiff"
	// OUT - the folder where recordings should be saved
	OUT = "/perm/recordings/"
)

func main() {
	// Logging
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.Info("starting weddingphone...")

	// concurrency handling
	main := context.Background()

	// Initialize audio system
	portaudio.Initialize()

	// First intro sound
	// play.PlayAIFF(main, INTRO)

	// Figure out a recording file
	recordingFile := path.Join(OUT, fmt.Sprintf("%d.%s", time.Now().Unix(), "aiff"))

	// Record the message
	recordCtx, cancelRecord := context.WithCancel(main) // Create a new child context from main
	go record.RecordToFile(recordCtx, recordingFile)    // start recording in the background
	time.Sleep(10 * time.Second)                        // Wait 10s
	cancelRecord()                                      // Then stop the recording

	// Play the message
	time.Sleep(1 * time.Second)
	play.PlayAIFF(main, recordingFile)

	// Shutodwn rountine
	portaudio.Terminate()
}
