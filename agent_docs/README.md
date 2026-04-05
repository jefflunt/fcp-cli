# fcp-cli — Agent Documentation Index

**fcp-cli** is a CLI tool to automate video creation by combining Final Cut Pro (via FCPXML), ffmpeg, and Apple Compressor. It handles audio normalization, duration probing, FCPXML generation, and Compressor rendering.

## Quick Start

### Installation
You need Go 1.22+ and external dependencies (ffmpeg, ffprobe, Compressor).

```bash
go build -o fcp-cli main.go
```

### Usage
Create a `config.yaml` file (see `config.yaml.example`) and run:

```bash
./fcp-cli -config config.yaml
```

---

## How to Use This Documentation

This folder follows **Progressive Disclosure** principles. Start here, then read the detailed files relevant to your task.

| File | What it covers | Read when… |
|------|---------------|------------|
| **This file** | Repo overview, file map, key facts | Always — start here |
| [`architecture.md`](architecture.md) | Component overview, data flow | Understanding how the tool works |
| [`config.md`](config.md) | Configuration options | Configuring a new project |

---

## Repo at a Glance

```
fcp-cli/
├── audio/          ← Audio normalization logic
├── compressor/     ← Compressor rendering interface
├── config/         ← Configuration parsing
├── fcpxml/         ← FCPXML generation logic
├── probe/          ← Media file duration probing
├── validator/      ← Asset and dependency validation
├── agent_docs/     ← Documentation
├── main.go         ← Entry point & pipeline orchestration
└── config.yaml.example
```
