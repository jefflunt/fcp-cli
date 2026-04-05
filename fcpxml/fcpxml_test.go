package fcpxml_test

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jefflunt/fcp-cli/fcpxml"
)

func defaultParams(dir string) fcpxml.Params {
	return fcpxml.Params{
		ProjectName:     "TestProject",
		LibraryPath:     "/Users/jeff/Movies/Main.fclib",
		TitleCardPath:   filepath.Join(dir, "title.png"),
		IntroVOPath:     filepath.Join(dir, "intro.wav"),
		DevlogPath:      filepath.Join(dir, "session.mov"),
		NormAudioPath:   filepath.Join(dir, "session_norm.wav"),
		FPS:             30,
		IntroDurTicks:   300,  // 10 seconds at 30fps
		DevlogDurTicks:  1800, // 60 seconds at 30fps
		TransitionTicks: 120,  // 4 seconds at 30fps
	}
}

func TestGenerate_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.fcpxml")

	if err := fcpxml.Generate(defaultParams(dir), outPath); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Fatal("expected output file to be created")
	}
}

func TestGenerate_ValidXML(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.fcpxml")

	if err := fcpxml.Generate(defaultParams(dir), outPath); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	var doc interface{}
	if err := xml.Unmarshal(data, &doc); err != nil {
		t.Errorf("output is not valid XML: %v", err)
	}
}

func TestGenerate_ContainsExpectedContent(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.fcpxml")

	params := defaultParams(dir)
	if err := fcpxml.Generate(params, outPath); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	content := string(data)

	checks := []struct {
		desc    string
		snippet string
	}{
		{"FCPXML version", `version="1.10"`},
		{"project name", "TestProject"},
		{"format frameDuration", `frameDuration="1/30s"`},
		{"4K width", `width="3840"`},
		{"4K height", `height="2160"`},
		{"library path URI", "file:///Users/jeff/Movies/Main.fclib"},
		{"title card resource", `name="title_card"`},
		{"intro VO resource", `name="intro_vo"`},
		{"devlog video resource", `name="devlog_video"`},
		{"norm audio resource", `name="norm_audio"`},
		{"cross dissolve transition", "Cross Dissolve"},
		{"dialogue role", `role="dialogue"`},
		{"intro duration", "300/30s"},
		{"devlog duration", "1800/30s"},
		{"transition duration", "120/30s"},
	}

	for _, c := range checks {
		if !strings.Contains(content, c.snippet) {
			t.Errorf("missing %s: expected to find %q in output", c.desc, c.snippet)
		}
	}
}

func TestGenerate_CreatesOutputDir(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "subdir", "nested", "out.fcpxml")

	if err := fcpxml.Generate(defaultParams(dir), outPath); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Fatal("expected output file to be created in nested directory")
	}
}

func TestGenerate_FileURIPaths(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.fcpxml")

	if err := fcpxml.Generate(defaultParams(dir), outPath); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	data, _ := os.ReadFile(outPath)
	content := string(data)

	// All asset src attributes must use file:// URIs.
	if !strings.Contains(content, "file://") {
		t.Error("expected file:// URIs in asset src attributes")
	}
}
