package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AssetDir           string `mapstructure:"asset_dir"`
	ControlFlowPath    string `mapstructure:"control_flow_path"`
	RuntimeContextPath string `mapstructure:"runtime_context_path"`
	OutputPath         string `mapstructure:"output_path"`
	FontDir            string `mapstructure:"font_dir"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("asset_dir", "assets/")
	viper.SetDefault("control_flow_path", "flows/section_oriented_control_flow.json")
	viper.SetDefault("runtime_context_path", "contexts/runtime_context.json")
	viper.SetDefault("output_path", "output/document.pdf")
	viper.SetDefault("font_dir", "assets/fonts")

	viper.SetEnvPrefix("GOPDF")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is present, use it
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Also support .env files if present
	viper.SetConfigFile(".env")
	if err := viper.MergeInConfig(); err != nil {
		// Ignore if .env is missing, we already have defaults and maybe config.yaml
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
