package sound

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/go-ini/ini"
)

const (
	OSReleaseFile = "/etc/os-release"
)

// AlsaSound implements SoundSystem using alsa-utils
type AlsaSound struct{}

func NewAlsaSound() *AlsaSound {
	// make sure that AlsaSound works properly by making sure alsa-utils is installed
	log.Print("Installing alsa-utils for correct distro and ensuring binaries are in PATH")
	// osInfo, err := readOSRelease(OSReleaseFile)
	// if err != nil {
	// 	return err
	// }
	// osRelease := osInfo["ID"]

	// we don't check for installed packages, we just reinstall
	// switch osRelease {
	// case "arch":
	// 	cmd := exec.Command("sudo", "pacman", "-S", "--no-confirm", "alsa-utils")
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		return err
	// 	}
	// case "debian":
	// 	cmd := exec.Command("sudo", "apt", "install", "-y", "alsa-utils")
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return &AlsaSound{}
}

func (as *AlsaSound) PlayWAV(ctx context.Context, filePath string) error {
	log.Printf("Playing file %s ...", filePath)
	return nil
}

func (as *AlsaSound) RecordToFile(ctx context.Context, filePath string) error {
	cmd := exec.Command("arecord", filePath)
	err := cmd.Start()
	if err != nil {
		log.Errorf("recording failed to start: %q", err)
	} else {
		log.Print("recording has been started")
	}

	<-ctx.Done()
	// stop the recording
	syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("recording failed: %q", err)
	}
	return nil
}

// Source: https://stackoverflow.com/questions/42386788/how-can-i-detect-linux-distribution-within-golang-program#42387003
func readOSRelease(configfile string) (map[string]string, error) {
	cfg, err := ini.Load(configfile)
	if err != nil {
		return nil, fmt.Errorf("fail to read file: %q", err)
	}

	configParams := make(map[string]string)
	configParams["ID"] = cfg.Section("").Key("ID").String()

	return configParams, nil
}
