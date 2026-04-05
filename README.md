# fcp-cli
Final Cut Pro + FCPXML + ffmpeg + Compressor to automate the creation of videos for upload.

## Quick Start

### Installation
You need Go 1.22+ and external dependencies (ffmpeg, ffprobe, Compressor).

```bash
go build -o fcp-cli main.go
```

### Usage
1. Copy `config.yaml.example` to `config.yaml` and update with your settings:

```yaml
project_name: "Devlog_042"
library_path: "/Users/jeff/Movies/Main.fclib"
output_dir: "./exports"
assets:
  intro_vo: "./assets/intro_voiceover.wav"
  title_card: "./assets/title_card.png"
  devlog_video: "./raw/session_capture.mov"
settings:
  fps: 30
  target_lufs: -14.0
  transition_duration_seconds: 4
  compressor_spec_path: "Built-In/Apple Devices/4K (H.264)"
```

2. Run the tool:

```bash
./fcp-cli -config config.yaml
```

See [agent_docs/](agent_docs/README.md) for detailed architectural documentation.
