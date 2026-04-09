package render

import (
	"fmt"
	"os/exec"
)

// Render combines the video and normalized audio into the final output file
// using ffmpeg.
func Render(videoPath, audioPath, outputPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", videoPath,
		"-i", audioPath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-b:a", "192k",
		"-map", "0:v:0",
		"-map", "1:a:0",
		"-shortest",
		outputPath,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg render failed: %w\n%s", err, out)
	}

	return nil
}
