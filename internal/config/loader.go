package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	"github.com/cruffinoni/rimworld-editor/pkg/runtimepath"
)

const envPrefix = "RIMEDIT"

func configureDefaults(v *viper.Viper) {
	defaults := map[string]any{
		"logging.level":  "info",
		"logging.format": "text",
		"logging.output": "stdout",
	}

	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}

func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return err
	}
	return nil
}

const defaultTagName = "mapstructure"

func Load() (*Config, error) {
	v := viper.New()
	configureDefaults(v)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/therapon/")
	v.AddConfigPath(runtimepath.ConfigDir())
	v.AddConfigPath("$HOME/.therapon/")
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}
	var cfg Config
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToURLHookFunc(),
		),
		Result:  &cfg,
		TagName: defaultTagName,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}
	if err = decoder.Decode(v.AllSettings()); err != nil {
		return nil, err
	}
	if err = validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
