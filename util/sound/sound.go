package sound

import "context"

// SoundSystem represents a way to play/record sound on a raspberry pi
type SoundSystem interface {
	// RecordToFile takes the given filepath, creates the file and write WAV comtabile audio to it, until context is cancelled
	RecordToFile(ctx context.Context, filePath string) error
	// PlayWAV playes the given file until context is canceled
	PlayWAV(ctx context.Context, filePath string) error
}
