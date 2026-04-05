# fcp-cli
Final Cut Pro + FCPXML + ffmpeg + Compressor to automate the creation of videos for upload.

## Quick Start

### Installation
You need Go 1.22+ and external dependencies (ffmpeg, ffprobe, Compressor).

```bash
go build -o fcp-cli main.go
```

### Usage
1. Copy `config.yaml.example` to `config.yaml` and update with your settings.
2. Run the tool:

```bash
./fcp-cli -config config.yaml
```

See [agent_docs/](agent_docs/README.md) for detailed architectural documentation.
