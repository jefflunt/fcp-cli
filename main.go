package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
	"path/filepath"

	"github.com/jefflunt/fcp-cli/audio"
	"github.com/jefflunt/fcp-cli/config"
	"github.com/jefflunt/fcp-cli/probe"
	"github.com/jefflunt/fcp-cli/render"
	"github.com/jefflunt/fcp-cli/validator"
)

func main() {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "config.yaml", "path to the YAML configuration file")
	keepWAV := pflag.Bool("keep-wav", false, "keep intermediate normalized WAV file after rendering")
	pflag.Parse()

	log.SetFlags(0)
	log.SetPrefix("[fcp-cli] ")

	if err := run(configPath, *keepWAV); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run(configPath string, keepWAV bool) error {
	// ── Step 1: Load & validate configuration ────────────────────────────────
	log.Printf("loading config from %s", configPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// ── Step 1a: Validate external dependencies ───────────────────────────────
	log.Println("checking dependencies (ffmpeg, ffprobe)...")
	if err := validator.CheckDependencies(); err != nil {
		return fmt.Errorf("dependency check failed: %w", err)
	}

	// ── Step 1b: Validate asset files ─────────────────────────────────────────
	log.Println("validating configuration paths...")
	var introVOPath string
	if s, ok := cfg.Assets.IntroVO.(string); ok {
		introVOPath = s
	}

	if err := validator.CheckPaths(
		cfg.LibraryPath,
		cfg.OutputDir,
		cfg.Assets.TitleCard,
		cfg.Assets.DevlogVideo,
		introVOPath,
	); err != nil {
		return fmt.Errorf("config path check failed: %w", err)
	}
	log.Println("all configuration paths valid")

	log.Println("checking asset files...")
	if err := validator.CheckAssets(
		introVOPath,
		cfg.Assets.TitleCard,
		cfg.Assets.DevlogVideo,
	); err != nil {
		return fmt.Errorf("asset check failed: %w", err)
	}

	// Resolve all asset paths to absolute paths so FCPXML URIs are correct.
	var introVOAbs string
	if introVOPath != "" {
		var err error
		introVOAbs, err = filepath.Abs(introVOPath)
		if err != nil {
			return fmt.Errorf("resolving intro_vo path: %w", err)
		}
	}
	titleCardAbs, err := filepath.Abs(cfg.Assets.TitleCard)
	if err != nil {
		return fmt.Errorf("resolving title_card path: %w", err)
	}
	_ = titleCardAbs // Keep variable to minimize diff, or remove it

	devlogVideoAbs, err := filepath.Abs(cfg.Assets.DevlogVideo)
	if err != nil {
		return fmt.Errorf("resolving devlog_video path: %w", err)
	}
	outputDirAbs, err := filepath.Abs(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("resolving output_dir path: %w", err)
	}

	// ── Step 2: Audio normalization ───────────────────────────────────────────
	normWAVPath := filepath.Join(outputDirAbs, cfg.ProjectName+"_norm.wav")
	log.Printf("normalizing audio → %s", normWAVPath)

	if err := os.MkdirAll(outputDirAbs, 0o755); err != nil {
		return fmt.Errorf("creating output directory %s: %w", outputDirAbs, err)
	}

	if err := audio.Normalize(devlogVideoAbs, cfg.Settings.TargetLUFS, normWAVPath); err != nil {
		return fmt.Errorf("audio normalization: %w", err)
	}

	// ── Step 3: Probe durations ───────────────────────────────────────────────
	var introDurSec float64
	var introDurTicks int64
	var transitionTicks int64
	if introVOAbs != "" {
		log.Printf("probing duration of intro_vo: %s", introVOAbs)
		var err error
		introDurSec, err = probe.Duration(introVOAbs)
		if err != nil {
			return fmt.Errorf("probing intro_vo duration: %w", err)
		}
		introDurTicks = probe.SecondsToTicks(introDurSec, cfg.Settings.FPS)
		transitionTicks = int64(float64(cfg.Settings.TransitionDurationSeconds) * float64(cfg.Settings.FPS))
	}

	log.Printf("probing duration of devlog_video: %s", devlogVideoAbs)
	devlogDurSec, err := probe.Duration(devlogVideoAbs)
	if err != nil {
		return fmt.Errorf("probing devlog_video duration: %w", err)
	}

	fps := cfg.Settings.FPS
	devlogDurTicks := probe.SecondsToTicks(devlogDurSec, fps)

	log.Printf("devlog duration: %.3fs (%d ticks @ %d fps)", devlogDurSec, devlogDurTicks, fps)
	if introVOAbs != "" {
		log.Printf("intro duration:  %.3fs (%d ticks @ %d fps)", introDurSec, introDurTicks, fps)
		log.Printf("transition:      %d ticks @ %d fps", transitionTicks, fps)
	}

	// ── Step 4: Render video ──────────────────────────────────────────────────
	outputPath := filepath.Join(outputDirAbs, "final_render.mp4")
	log.Printf("rendering final video → %s", outputPath)

	if err := render.Render(devlogVideoAbs, normWAVPath, outputPath); err != nil {
		return fmt.Errorf("rendering video: %w", err)
	}

	log.Printf("done! output written to %s", outputPath)

	// ── Cleanup ───────────────────────────────────────────────────────────────
	if !keepWAV {
		log.Printf("removing intermediate WAV: %s", normWAVPath)
		if err := os.Remove(normWAVPath); err != nil {
			// Non-fatal: log but do not fail.
			log.Printf("warning: could not remove intermediate WAV %s: %v", normWAVPath, err)
		}
	}

	return nil
}
