package validator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jefflunt/fcp-cli/validator"
)

func TestCheckAssets_AllPresent(t *testing.T) {
	dir := t.TempDir()

	introVO := filepath.Join(dir, "intro.wav")
	titleCard := filepath.Join(dir, "title.png")
	devlogVideo := filepath.Join(dir, "session.mov")

	for _, p := range []string{introVO, titleCard, devlogVideo} {
		if err := os.WriteFile(p, []byte{}, 0o644); err != nil {
			t.Fatalf("creating fixture %s: %v", p, err)
		}
	}

	if err := validator.CheckAssets(introVO, titleCard, devlogVideo); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestCheckAssets_MissingFile(t *testing.T) {
	dir := t.TempDir()

	introVO := filepath.Join(dir, "intro.wav")
	titleCard := filepath.Join(dir, "title.png")
	devlogVideo := filepath.Join(dir, "missing.mov") // does not exist

	for _, p := range []string{introVO, titleCard} {
		if err := os.WriteFile(p, []byte{}, 0o644); err != nil {
			t.Fatalf("creating fixture %s: %v", p, err)
		}
	}

	if err := validator.CheckAssets(introVO, titleCard, devlogVideo); err == nil {
		t.Error("expected error for missing asset, got nil")
	}
}

func TestCheckPaths_AllPresent(t *testing.T) {
	dir := t.TempDir()

	p1 := filepath.Join(dir, "lib.fclib")
	p2 := filepath.Join(dir, "output")
	p3 := filepath.Join(dir, "title.png")
	p4 := filepath.Join(dir, "video.mov")
	p5 := filepath.Join(dir, "compressor.app")

	for _, p := range []string{p1, p2, p3, p4, p5} {
		if err := os.WriteFile(p, []byte{}, 0o644); err != nil {
			t.Fatalf("creating fixture %s: %v", p, err)
		}
	}

	if err := validator.CheckPaths(p1, p2, p3, p4, ""); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
