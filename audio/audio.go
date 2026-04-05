package audio

import (
	"fmt"
	"os/exec"
)

// Normalize runs ffmpeg with the loudnorm filter on the devlog video, writing
// a normalized PCM WAV file to outPath. targetLUFS is the integrated loudness
// target (e.g. -14.0).
func Normalize(devlogVideo string, targetLUFS float64, outPath string) error {
	filter := fmt.Sprintf("loudnorm=I=%.1f:TP=-1.5:LRA=11", targetLUFS)

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", devlogVideo,
		"-af", filter,
		"-vn",
		"-c:a", "pcm_s16le",
		outPath,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg loudnorm failed: %w\n%s", err, out)
	}

	return nil
}
