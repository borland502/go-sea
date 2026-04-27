package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/borland502/go-sea/internal/appconfig"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName                  = "go-sea"
	binaryName               = "go-sea"
	envPrefix                = "GOSEA"
	configBaseName           = "config"
	configExtension          = "toml"
	exampleSectionKey        = "example"
	outputSectionKey         = "output"
	paletteSectionKey        = "palette"
	paletteName              = "Monokai Spectrumish"
	skipConfigLoadAnnotation = "skip-config-load"
)

type Config = appconfig.Config

type ExampleConfig = appconfig.ExampleConfig

type OutputConfig = appconfig.OutputConfig

type PaletteConfig = appconfig.PaletteConfig

type PaletteEntry = appconfig.PaletteEntry

var errConfigFileNotFound = errors.New("config file not found")

var (
	cfgFile          string
	loadedConfigPath string
	cfg              = defaultConfig()
	searchConfigFile = searchConfigFileInXDGPaths
)

func defaultConfig() Config {
	return appconfig.Default()
}

func configFileName() string {
	return configBaseName + "." + configExtension
}

func configRelativePath() string {
	return filepath.Join(appName, configFileName())
}

func defaultUserConfigFilePath() string {
	return filepath.Join(xdg.ConfigHome, configRelativePath())
}

func configSourceDescription() string {
	if loadedConfigPath != "" {
		return loadedConfigPath
	}

	return "defaults, environment variables, and flags"
}

func shouldSkipConfigLoad(cmd *cobra.Command) bool {
	return cmd.Annotations[skipConfigLoadAnnotation] == "true"
}

func loadConfig(cmd *cobra.Command) error {
	loader, err := newConfigLoader(cmd)
	if err != nil {
		return err
	}

	loadedConfigPath = ""

	shouldRead, err := configureConfigFile(loader)
	if err != nil {
		return err
	}

	if shouldRead {
		if err := loader.ReadInConfig(); err != nil {
			return fmt.Errorf("read config: %w", err)
		}
		loadedConfigPath = loader.ConfigFileUsed()
	}

	cfg = defaultConfig()
	if err := loader.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("decode config: %w", err)
	}

	return nil
}

func newConfigLoader(cmd *cobra.Command) (*viper.Viper, error) {
	defaults := defaultConfig()
	loader := viper.New()

	loader.SetEnvPrefix(envPrefix)
	loader.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	loader.AutomaticEnv()

	loader.SetDefault(exampleSectionKey+".greeting", defaults.Example.Greeting)
	loader.SetDefault(exampleSectionKey+".name", defaults.Example.Name)
	loader.SetDefault(exampleSectionKey+".location", defaults.Example.Location)
	loader.SetDefault(exampleSectionKey+".favorite_tide", defaults.Example.FavoriteTide)
	loader.SetDefault(outputSectionKey+".color", defaults.Output.Color)
	loader.SetDefault(outputSectionKey+".show_paths", defaults.Output.ShowPaths)

	for _, entry := range defaults.Palette.Entries() {
		loader.SetDefault(paletteSectionKey+"."+entry.Key, entry.Hex)
	}

	for _, key := range []string{
		exampleSectionKey + ".greeting",
		exampleSectionKey + ".name",
		exampleSectionKey + ".location",
		exampleSectionKey + ".favorite_tide",
		outputSectionKey + ".color",
		outputSectionKey + ".show_paths",
	} {
		if err := loader.BindEnv(key); err != nil {
			return nil, fmt.Errorf("bind env %s: %w", key, err)
		}
	}

	for _, entry := range defaults.Palette.Entries() {
		key := paletteSectionKey + "." + entry.Key
		if err := loader.BindEnv(key); err != nil {
			return nil, fmt.Errorf("bind env %s: %w", key, err)
		}
	}

	if flag := cmd.Flags().Lookup("greeting"); flag != nil {
		if err := loader.BindPFlag(exampleSectionKey+".greeting", flag); err != nil {
			return nil, fmt.Errorf("bind greeting flag: %w", err)
		}
	}
	if flag := cmd.Flags().Lookup("name"); flag != nil {
		if err := loader.BindPFlag(exampleSectionKey+".name", flag); err != nil {
			return nil, fmt.Errorf("bind name flag: %w", err)
		}
	}
	if flag := cmd.Flags().Lookup("location"); flag != nil {
		if err := loader.BindPFlag(exampleSectionKey+".location", flag); err != nil {
			return nil, fmt.Errorf("bind location flag: %w", err)
		}
	}
	if flag := cmd.Flags().Lookup("favorite-tide"); flag != nil {
		if err := loader.BindPFlag(exampleSectionKey+".favorite_tide", flag); err != nil {
			return nil, fmt.Errorf("bind favorite-tide flag: %w", err)
		}
	}
	if flag := cmd.Flags().Lookup("color"); flag != nil {
		if err := loader.BindPFlag(outputSectionKey+".color", flag); err != nil {
			return nil, fmt.Errorf("bind color flag: %w", err)
		}
	}
	if flag := cmd.Flags().Lookup("show-paths"); flag != nil {
		if err := loader.BindPFlag(outputSectionKey+".show_paths", flag); err != nil {
			return nil, fmt.Errorf("bind show-paths flag: %w", err)
		}
	}

	return loader, nil
}

func configureConfigFile(loader *viper.Viper) (bool, error) {
	if cfgFile != "" {
		loader.SetConfigFile(cfgFile)
		return true, nil
	}

	configPath, err := searchConfigFile(configRelativePath())
	if err != nil {
		if errors.Is(err, errConfigFileNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("search config file: %w", err)
	}

	loader.SetConfigFile(configPath)
	return true, nil
}

func searchConfigFileInXDGPaths(relPath string) (string, error) {
	searchPaths := append([]string{xdg.ConfigHome}, xdg.ConfigDirs...)
	searchedPaths := make([]string, 0, len(searchPaths))

	for _, basePath := range searchPaths {
		candidatePath := filepath.Join(basePath, relPath)
		info, err := os.Stat(candidatePath)
		if err == nil {
			if info.IsDir() {
				return "", fmt.Errorf("config path is a directory: %s", candidatePath)
			}

			return candidatePath, nil
		}
		if errors.Is(err, os.ErrNotExist) {
			searchedPaths = append(searchedPaths, filepath.Dir(candidatePath))
			continue
		}

		return "", fmt.Errorf("stat config path %s: %w", candidatePath, err)
	}

	return "", fmt.Errorf("%w: %s", errConfigFileNotFound, strings.Join(searchedPaths, ", "))
}

func writableConfigFilePath() (string, error) {
	if cfgFile != "" {
		if err := os.MkdirAll(filepath.Dir(cfgFile), 0o700); err != nil {
			return "", fmt.Errorf("create config directory: %w", err)
		}
		return cfgFile, nil
	}

	configPath, err := xdg.ConfigFile(configRelativePath())
	if err != nil {
		return "", fmt.Errorf("resolve XDG config path: %w", err)
	}

	return configPath, nil
}

func renderConfig(config Config) string {
	return appconfig.Render(config, binaryName)
}
