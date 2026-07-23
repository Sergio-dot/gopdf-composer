// Package config provides 12-factor configuration via Viper.
// Settings are loaded from YAML files, .env files, and GOPDF_ prefixed
// environment variables, with sensible defaults for all values.
package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the PDF generation engine.
type Config struct {
	AssetDir           string                       `mapstructure:"asset_dir"`
	ControlFlowPath    string                       `mapstructure:"control_flow_path"`
	RuntimeContextPath string                       `mapstructure:"runtime_context_path"`
	OutputPath         string                       `mapstructure:"output_path"`
	FontDir            string                       `mapstructure:"font_dir"`
	DefaultFont        string                       `mapstructure:"default_font"`
	FontFiles          map[string]map[string]string `mapstructure:"font_files"`
}

// LoadConfig reads configuration from YAML, .env, and GOPDF_ environment
// variables. Returns a Config with sensible defaults when no files are present.
func LoadConfig() (*Config, error) {
	viper.SetDefault("asset_dir", "assets/")
	viper.SetDefault("control_flow_path", "flows/flow.json")
	viper.SetDefault("runtime_context_path", "contexts/context.json")
	viper.SetDefault("output_path", "output/document.pdf")
	viper.SetDefault("font_dir", "assets/fonts")
	viper.SetDefault("default_font", "Arial")

	viper.SetEnvPrefix("GOPDF")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err != nil {
		// .env is optional, defaults and config.yaml already loaded
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
