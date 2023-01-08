package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/the-technat/weddingphone/util/device"
	"github.com/the-technat/weddingphone/util/sound"
)

const ()

var (
	SAVE_PATH     = "./dist/recordings"
	mainCtx, stop = context.WithCancel(context.Background())
	recordChannel = make(chan bool)
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// make sure output directory is set
	savePath, found := os.LookupEnv("SAVE_PATH")
	if found {
		// check if dir also exists
		_, err := os.Stat(savePath)
		if err != nil {
			os.MkdirAll(savePath, os.FileMode(0751))
			log.Warningf("%s doesn't existed, created it")
			SAVE_PATH = savePath
		}
	}
	if !found {
		err := os.MkdirAll(SAVE_PATH, os.FileMode(0751))
		if err != nil {
			log.Fatal(err)
		}
		log.Warningf("SAVE_PATH not set, using default of %s", SAVE_PATH)
	}
}

func main() {
	// Device initalization according to ../docs/hardware.md
	d := device.New(
		map[string]device.LED{
			"status": {
				Color: "green",
				Pin:   2,
			},
			"recording": {
				Color: "orange",
				Pin:   3,
			},
		},
		map[string]device.Button{
			"recording": {
				Pin:                     4,
				EnableFallEdgeDetection: true,
			},
			"play": {
				Pin:                     27,
				EnableFallEdgeDetection: true,
			},
		},
	)
	defer d.Close()

	// Sound initalization
	ss := sound.NewAlsaSound()

	log.Print("started weddingphone")
	defer stop()
	d.BlinkDuration(mainCtx, "status", time.Second/5, 2*time.Second)
	d.LEDs["status"].Toggle()

	// continuously monitor the button for a trigger
	// notify the recordChannel in case of a buttonPress
	go d.ButtonEvents(mainCtx, "recording", recordChannel, time.Second)

	// continuously check if we got notifed about a button press
	// if so, record until the next button press
	go recordRoutine(mainCtx, recordChannel, d, ss)

	// Shutdown routine
	done := make(chan struct{})
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		// wait for someone to request shutdown (e.g ctrl+c or systemd)
		<-shutdown
		// if shutdown is requested, close the main ctx which should cause all sub-process (which derive from mainCtx) to stop
		stop()
		// and notify the main programm to stop as well
		close(done)
	}()

	<-done
	log.Print("stopped weddingphone")
}

func recordRoutine(ctx context.Context, notify chan bool, d *device.Device, ss sound.SoundSystem) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-notify:
			// start blinking
			blinkerCtx, stopBlinker := context.WithCancel(ctx)
			go d.Blink(blinkerCtx, "recording", time.Second)

			// start recording
			recordingFile := path.Join(SAVE_PATH, fmt.Sprintf("%d.%s", time.Now().Unix(), "wav"))
			recordCtx, stopRecording := context.WithCancel(ctx)
			go ss.RecordToFile(recordCtx, recordingFile)

			<-notify
			stopRecording()
			stopBlinker()
			// d.LEDs["recording"].PullDown()
			log.Print("recording stopped")
		}
	}
}
