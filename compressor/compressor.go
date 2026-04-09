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

	// Use a simpler setting path if possible, but let's test with the one provided first
	// but use -batchname to avoid the batch identifier problem.
	cmd := exec.Command(
		binary,
		"-batchname", "fcp-cli-render",
		"-jobpath", fcpxmlPath,
		"-settingpath", "/tmp/setting.compressorsetting",
		"-locationpath", outputPath,
		"-outputformat", "xml",
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Compressor render failed: %w\n%s", err, out)
	}

	return nil
}
