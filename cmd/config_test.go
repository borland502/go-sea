package cmd

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigureConfigFilePropagatesSearchErrors(t *testing.T) {
	originalCfgFile := cfgFile
	originalSearch := searchConfigFile
	t.Cleanup(func() {
		cfgFile = originalCfgFile
		searchConfigFile = originalSearch
	})

	cfgFile = ""
	expectedErr := errors.New("boom")
	searchConfigFile = func(relPath string) (string, error) {
		if relPath != filepath.Join(appName, configFileName()) {
			t.Fatalf("unexpected relative path %q", relPath)
		}

		return "", expectedErr
	}

	_, err := configureConfigFile(viper.New())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected search error %v, got %v", expectedErr, err)
	}
}

func TestConfigureConfigFileIgnoresConfigNotFound(t *testing.T) {
	originalCfgFile := cfgFile
	originalSearch := searchConfigFile
	t.Cleanup(func() {
		cfgFile = originalCfgFile
		searchConfigFile = originalSearch
	})

	cfgFile = ""
	searchConfigFile = func(relPath string) (string, error) {
		if relPath != filepath.Join(appName, configFileName()) {
			t.Fatalf("unexpected relative path %q", relPath)
		}

		return "", errConfigFileNotFound
	}

	shouldRead, err := configureConfigFile(viper.New())
	if err != nil {
		t.Fatalf("expected no error for missing config, got %v", err)
	}
	if shouldRead {
		t.Fatal("expected missing config to skip reading")
	}
}

func TestConfigureConfigFileUsesExplicitPath(t *testing.T) {
	originalCfgFile := cfgFile
	originalSearch := searchConfigFile
	t.Cleanup(func() {
		cfgFile = originalCfgFile
		searchConfigFile = originalSearch
	})

	configPath := filepath.Join(t.TempDir(), "config.toml")
	cfgFile = configPath
	searchCalled := false
	searchConfigFile = func(string) (string, error) {
		searchCalled = true
		return "", nil
	}

	shouldRead, err := configureConfigFile(viper.New())
	if err != nil {
		t.Fatalf("expected no error for explicit config path, got %v", err)
	}
	if !shouldRead {
		t.Fatal("expected explicit config path to be read")
	}
	if searchCalled {
		t.Fatal("expected explicit config path to bypass XDG search")
	}
}
