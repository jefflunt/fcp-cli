package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// expandTilde replaces the leading ~ with the user's home directory.
func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	if len(path) == 1 {
		return usr.HomeDir, nil
	}
	return filepath.Join(usr.HomeDir, path[2:]), nil
}

// Assets holds the file paths for all media assets used in the project.
type Assets struct {
	IntroVO     any    `yaml:"intro_vo"`
	TitleCard   string `yaml:"title_card"`
	DevlogVideo string `yaml:"devlog_video"`
}

// Settings holds the project's technical parameters.
type Settings struct {
	FPS                       int     `yaml:"fps"`
	TargetLUFS                float64 `yaml:"target_lufs"`
	TransitionDurationSeconds int     `yaml:"transition_duration_seconds"`
	CompressorSpecPath        string  `yaml:"compressor_spec_path"`
}

// Config is the top-level configuration loaded from a YAML file.
type Config struct {
	ProjectName string   `yaml:"project_name"`
	LibraryPath string   `yaml:"library_path"`
	OutputDir   string   `yaml:"output_dir"`
	Assets      Assets   `yaml:"assets"`
	Settings    Settings `yaml:"settings"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	// Expand tildes in paths
	cfg.LibraryPath, err = expandTilde(cfg.LibraryPath)
	if err != nil {
		return nil, fmt.Errorf("expanding library_path %q: %w", cfg.LibraryPath, err)
	}
	cfg.OutputDir, err = expandTilde(cfg.OutputDir)
	if err != nil {
		return nil, fmt.Errorf("expanding output_dir %q: %w", cfg.OutputDir, err)
	}
	if introVO, ok := cfg.Assets.IntroVO.(string); ok && introVO != "" {
		var expanded string
		expanded, err = expandTilde(introVO)
		if err != nil {
			return nil, fmt.Errorf("expanding intro_vo %q: %w", introVO, err)
		}
		cfg.Assets.IntroVO = expanded
	}
	cfg.Assets.TitleCard, err = expandTilde(cfg.Assets.TitleCard)
	if err != nil {
		return nil, fmt.Errorf("expanding title_card %q: %w", cfg.Assets.TitleCard, err)
	}
	cfg.Assets.DevlogVideo, err = expandTilde(cfg.Assets.DevlogVideo)
	if err != nil {
		return nil, fmt.Errorf("expanding devlog_video %q: %w", cfg.Assets.DevlogVideo, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.ProjectName == "" {
		return fmt.Errorf("project_name is required")
	}
	if c.LibraryPath == "" {
		return fmt.Errorf("library_path is required")
	}
	if c.OutputDir == "" {
		return fmt.Errorf("output_dir is required")
	}
	if val, ok := c.Assets.IntroVO.(string); ok {
		if val == "" {
			return fmt.Errorf("assets.intro_vo is required")
		}
	} else if val, ok := c.Assets.IntroVO.(bool); ok {
		if val {
			return fmt.Errorf("assets.intro_vo must be a string path or false")
		}
	} else {
		return fmt.Errorf("assets.intro_vo must be a string path or false")
	}
	if c.Assets.TitleCard == "" {
		return fmt.Errorf("assets.title_card is required")
	}
	if c.Assets.DevlogVideo == "" {
		return fmt.Errorf("assets.devlog_video is required")
	}
	if c.Settings.FPS <= 0 {
		return fmt.Errorf("settings.fps must be greater than zero")
	}
	if c.Settings.CompressorSpecPath == "" {
		return fmt.Errorf("settings.compressor_spec_path is required")
	}
	return nil
}
