package probe

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Duration returns the duration of a media file in seconds using ffprobe.
func Duration(filePath string) (float64, error) {
	out, err := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	).Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed for %q: %w", filePath, err)
	}

	durStr := strings.TrimSpace(string(out))
	dur, err := strconv.ParseFloat(durStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing duration %q from ffprobe: %w", durStr, err)
	}

	return dur, nil
}

// SecondsToTicks converts a duration in seconds to FCP timeline ticks at the
// given frames-per-second rate (e.g. 30). The returned value is the integer
// number of frames, ready to be formatted as "[n]/[fps]s" in FCPXML.
func SecondsToTicks(seconds float64, fps int) int64 {
	return int64(seconds * float64(fps))
}
