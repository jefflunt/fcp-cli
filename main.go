package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
	"path/filepath"

	"github.com/jefflunt/fcp-cli/audio"
	"github.com/jefflunt/fcp-cli/compressor"
	"github.com/jefflunt/fcp-cli/config"
	"github.com/jefflunt/fcp-cli/fcpxml"
	"github.com/jefflunt/fcp-cli/probe"
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
	log.Println("checking dependencies (ffmpeg, ffprobe, Compressor)...")
	if err := validator.CheckDependencies(); err != nil {
		return fmt.Errorf("dependency check failed: %w", err)
	}

	// ── Step 1b: Validate asset files ─────────────────────────────────────────
	log.Println("checking asset files...")
	if err := validator.CheckAssets(
		cfg.Assets.IntroVO,
		cfg.Assets.TitleCard,
		cfg.Assets.DevlogVideo,
	); err != nil {
		return fmt.Errorf("asset check failed: %w", err)
	}

	// Resolve all asset paths to absolute paths so FCPXML URIs are correct.
	introVOAbs, err := filepath.Abs(cfg.Assets.IntroVO)
	if err != nil {
		return fmt.Errorf("resolving intro_vo path: %w", err)
	}
	titleCardAbs, err := filepath.Abs(cfg.Assets.TitleCard)
	if err != nil {
		return fmt.Errorf("resolving title_card path: %w", err)
	}
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
	log.Printf("probing duration of intro_vo: %s", introVOAbs)
	introDurSec, err := probe.Duration(introVOAbs)
	if err != nil {
		return fmt.Errorf("probing intro_vo duration: %w", err)
	}

	log.Printf("probing duration of devlog_video: %s", devlogVideoAbs)
	devlogDurSec, err := probe.Duration(devlogVideoAbs)
	if err != nil {
		return fmt.Errorf("probing devlog_video duration: %w", err)
	}

	fps := cfg.Settings.FPS
	introDurTicks := probe.SecondsToTicks(introDurSec, fps)
	devlogDurTicks := probe.SecondsToTicks(devlogDurSec, fps)
	transitionTicks := int64(float64(cfg.Settings.TransitionDurationSeconds) * float64(fps))

	log.Printf("intro duration:  %.3fs (%d ticks @ %d fps)", introDurSec, introDurTicks, fps)
	log.Printf("devlog duration: %.3fs (%d ticks @ %d fps)", devlogDurSec, devlogDurTicks, fps)
	log.Printf("transition:      %d ticks @ %d fps", transitionTicks, fps)

	// ── Step 4: FCPXML generation ─────────────────────────────────────────────
	fcpxmlPath := filepath.Join(outputDirAbs, cfg.ProjectName+".fcpxml")
	log.Printf("generating FCPXML → %s", fcpxmlPath)

	if err := fcpxml.Generate(fcpxml.Params{
		ProjectName:     cfg.ProjectName,
		LibraryPath:     cfg.LibraryPath,
		TitleCardPath:   titleCardAbs,
		IntroVOPath:     introVOAbs,
		DevlogPath:      devlogVideoAbs,
		NormAudioPath:   normWAVPath,
		FPS:             fps,
		IntroDurTicks:   introDurTicks,
		DevlogDurTicks:  devlogDurTicks,
		TransitionTicks: transitionTicks,
	}, fcpxmlPath); err != nil {
		return fmt.Errorf("FCPXML generation: %w", err)
	}

	// ── Step 5: Compressor render ─────────────────────────────────────────────
	log.Printf("submitting to Compressor: %s", fcpxmlPath)
	if err := compressor.Render(fcpxmlPath, cfg.Settings.CompressorSpecPath, outputDirAbs); err != nil {
		return fmt.Errorf("Compressor render: %w", err)
	}

	log.Printf("done! output written to %s/final_render.mp4", outputDirAbs)

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
