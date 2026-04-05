package validator

import (
	"fmt"
	"os"
	"os/exec"
)

const compressorBinary = "/Applications/Compressor.app/Contents/MacOS/Compressor"

// CheckDependencies verifies that ffmpeg, ffprobe, and Apple Compressor are
// available on the host system.
func CheckDependencies() error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}
	if out, err := exec.Command("ffmpeg", "-version").CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg -version failed: %w\n%s", err, out)
	}

	if _, err := exec.LookPath("ffprobe"); err != nil {
		return fmt.Errorf("ffprobe not found in PATH: %w", err)
	}
	if out, err := exec.Command("ffprobe", "-version").CombinedOutput(); err != nil {
		return fmt.Errorf("ffprobe -version failed: %w\n%s", err, out)
	}

	if _, err := os.Stat(compressorBinary); os.IsNotExist(err) {
		return fmt.Errorf("Apple Compressor not found at %s", compressorBinary)
	}

	return nil
}

// CheckAssets verifies that all required asset files exist on disk.
func CheckAssets(introVO, titleCard, devlogVideo string) error {
	paths := map[string]string{
		"assets.intro_vo":     introVO,
		"assets.title_card":   titleCard,
		"assets.devlog_video": devlogVideo,
	}
	for field, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return fmt.Errorf("asset %s not found: %s", field, p)
		}
	}
	return nil
}
