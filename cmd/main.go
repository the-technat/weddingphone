package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gordonklaus/portaudio"
	log "github.com/sirupsen/logrus"
	"github.com/the-technat/weddingphone/util/play"
	"github.com/the-technat/weddingphone/util/record"
)

func main() {
	// Logging
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.Info("starting weddingphone...")

	// Read IN & OUT
	introPath := os.Getenv("INTRO_PATH")
	saveDir := os.Getenv("SAVE_PATH")
	if saveDir == "" || introPath == "" {
		log.Errorf("Env vars INTRO_PATH and/or SAVE_PATH not set")
	}

	// concurrency handling
	main := context.Background()

	// Initialize audio system
	portaudio.Initialize()

	// First intro sound
	play.PlayAIFF(main, introPath)

	// Figure out a recording file
	recordingFile := path.Join(saveDir, fmt.Sprintf("%d.%s", time.Now().Unix(), "aiff"))

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
