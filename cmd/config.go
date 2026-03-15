package cmd

import (
	"errors"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
)

// SubmoduleConfig holds per-submodule configuration from .wip.yml.
type SubmoduleConfig struct {
	OnWorktreeCreate []string `yaml:"on-worktree-create"`
}

// WipConfig is the top-level structure for .wip.yml.
type WipConfig struct {
	Submodules map[string]SubmoduleConfig `yaml:"submodules"`
}

// loadWipConfig reads .wip.yml from the current directory and parses it.
// Returns (nil, nil) if the file does not exist.
func loadWipConfig() (*WipConfig, error) {
	data, err := os.ReadFile(".wip.yml")
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
