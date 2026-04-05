# Architecture

The tool follows a linear, pipeline-based architecture:

1. **Config:** Loads and validates `config.yaml`.
2. **Validator:** Checks for external dependencies (ffmpeg, ffprobe, Compressor) and asset existence.
3. **Audio:** Normalizes audio from the input video to a target LUFS using ffmpeg.
4. **Probe:** Extracts durations from media files using ffprobe, converting them to FCPXML ticks.
5. **FCPXML:** Generates an FCPXML file referencing all assets and calculated timings.
6. **Compressor:** Submits the generated FCPXML file to Apple Compressor for final rendering.
