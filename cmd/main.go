package main

import (
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
	log "github.com/sirupsen/logrus"
	"github.com/the-technat/weddingphone/util/record"
)

const (
	// INTRO - the file that should be played as intro to the users
	INTRO = "/home/technat/code/weddingphone/assets/intro.aiff"
	OUT   = "assets/record.aiff"
)

func main() {
	// Logging
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.Info("starting weddingphone...")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	portaudio.Initialize()
	defer portaudio.Terminate()

	record.RecordToFile(OUT, sig)

}
