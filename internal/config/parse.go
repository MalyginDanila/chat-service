package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ParseAndValidate(filename string) (Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode file %q: %v", filename, err)
	}

	if envLevel, exists := os.LookupEnv("CHAT_LOG_LEVEL"); exists {
		cfg.Log.Level = envLevel
	}

	if err := validate.Struct(cfg); err != nil {
		return Config{}, fmt.Errorf("validate config: %v", err)
	}

	return cfg, nil
}
