package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jefflunt/fcp-cli/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

const validYAML = `
project_name: "TestProject"
library_path: "/tmp/test.fclib"
output_dir: "/tmp/exports"
assets:
  intro_vo: "./intro.wav"
  title_card: "./title.png"
  devlog_video: "./session.mov"
settings:
  fps: 30
  target_lufs: -14.0
  transition_duration_seconds: 4
  compressor_spec_path: "some/path"
final_render: "/tmp/final.mp4"
`

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, validYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.ProjectName != "TestProject" {
		t.Errorf("ProjectName = %q, want %q", cfg.ProjectName, "TestProject")
	}
	if cfg.Settings.FPS != 30 {
		t.Errorf("FPS = %d, want 30", cfg.Settings.FPS)
	}
	if cfg.Settings.TargetLUFS != -14.0 {
		t.Errorf("TargetLUFS = %f, want -14.0", cfg.Settings.TargetLUFS)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ":\tinvalid: yaml: [")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	cases := []struct {
		name    string
		content string
	}{
		{
			name: "missing project_name",
			content: `
library_path: "/tmp/test.fclib"
output_dir: "/tmp/out"
assets:
  intro_vo: "./a.wav"
  title_card: "./b.png"
  devlog_video: "./c.mov"
settings:
  fps: 30
  compressor_spec_path: "some/path"
`,
		},
		{
			name: "missing library_path",
			content: `
project_name: "P"
output_dir: "/tmp/out"
assets:
  intro_vo: "./a.wav"
  title_card: "./b.png"
  devlog_video: "./c.mov"
settings:
  fps: 30
  compressor_spec_path: "some/path"
`,
		},
		{
			name: "zero fps",
			content: `
project_name: "P"
library_path: "/tmp/test.fclib"
output_dir: "/tmp/out"
assets:
  intro_vo: "./a.wav"
  title_card: "./b.png"
  devlog_video: "./c.mov"
settings:
  fps: 0
  compressor_spec_path: "some/path"
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTemp(t, tc.content)
			_, err := config.Load(path)
			if err == nil {
				t.Fatalf("expected validation error for %q, got nil", tc.name)
			}
		})
	}
}
