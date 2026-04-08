package compressor

import (
	"fmt"
	"os/exec"
)

const binary = "/Applications/Compressor.app/Contents/MacOS/Compressor"

// Render submits an FCPXML file to Apple Compressor using the specified
// setting path and writes the rendered output to outDir/final_render.mp4.
func Render(fcpxmlPath, settingPath, outDir string) error {
	outputPath := outDir + "/final_render.mp4"

	// The Compressor CLI documentation indicates -locationpath is the correct flag.
	cmd := exec.Command(
		binary,
		"-jobpath", fcpxmlPath,
		"-settingpath", settingPath,
		"-locationpath", outputPath,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Compressor render failed: %w\n%s", err, out)
	}

	return nil
}
