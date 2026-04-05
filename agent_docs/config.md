# Configuration

The tool is configured via a YAML file (default `config.yaml`).

```yaml
projectName: "MyProject"
libraryPath: "/path/to/my.fcpbundle"
outputDir: "./output"
assets:
  introVO: "./assets/intro.wav"
  titleCard: "./assets/title.png"
  devlogVideo: "./assets/devlog.mp4"
settings:
  targetLUFS: -14.0
  fps: 30
  transitionDurationSeconds: 1.0
  compressorSpecPath: "/path/to/my.setting"
```
