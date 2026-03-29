package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"
)

// WipProject represents a located wip project — its root directory and parsed config.
type WipProject struct {
	Root   string
	Config *WipConfig
}

// FindWipProject locates the nearest .wip.yml by walking upward from the current
// working directory, bounded by the user's home directory. Returns an error if the
// cwd is outside the home directory, or if no .wip.yml is found within the tree.
func FindWipProject() (*WipProject, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not determine working directory: %w", err)
	}

	// Fail immediately if cwd is outside the home directory.
	if cwd != home && !strings.HasPrefix(cwd, home+string(filepath.Separator)) {
		return nil, errors.New("no .wip.yml found — run wip init first")
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, ".wip.yml")
		if _, err := os.Stat(candidate); err == nil {
			cfg, err := loadWipConfigFrom(dir)
			if err != nil {
				return nil, err
			}
			return &WipProject{Root: dir, Config: cfg}, nil
		}

		if dir == home {
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Filesystem root guard — shouldn't be reached given home ceiling above.
			break
		}
		dir = parent
	}

	return nil, errors.New("no .wip.yml found — run wip init first")
}

// RefConfig holds per-ref configuration from .wip.yml.
type RefConfig struct {
	URL              string   `yaml:"url"`
	Branch           string   `yaml:"branch"`
	OnWorktreeCreate []string `yaml:"on-worktree-create"`
	OnWorktreeLaunch []string `yaml:"on-worktree-launch"`
}

// WipConfig is the top-level structure for .wip.yml.
type WipConfig struct {
	Refs map[string]RefConfig `yaml:"refs"`
}

// loadWipConfigFrom reads .wip.yml from the given directory and parses it.
// Returns (nil, nil) if the file does not exist.
func loadWipConfigFrom(dir string) (*WipConfig, error) {
	data, err := os.ReadFile(filepath.Join(dir, ".wip.yml"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cfg WipConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// loadWipConfig reads .wip.yml from the current directory and parses it.
// Returns (nil, nil) if the file does not exist.
func loadWipConfig() (*WipConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return loadWipConfigFrom(cwd)
}

// requireWipConfig calls loadWipConfig and returns a user-facing error if the
// config file is absent.
func requireWipConfig() (*WipConfig, error) {
	cfg, err := loadWipConfig()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("no .wip.yml found — run wip init first")
	}
	return cfg, nil
}

// saveWipConfig marshals cfg to YAML and writes it to .wip.yml.
func saveWipConfig(cfg *WipConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(".wip.yml", data, 0644)
}

// stringList is a flag.Value implementation for repeatable string flags.
type stringList []string

func (s *stringList) String() string {
	return strings.Join([]string(*s), ",")
}

func (s *stringList) Set(val string) error {
	*s = append(*s, val)
	return nil
}
