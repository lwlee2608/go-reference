package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lwlee2608/adder"
	internalhttp "github.com/lwlee2608/go-reference/internal/api/http"
	"github.com/lwlee2608/go-reference/internal/db"
)

type Config struct {
	Log  LogConfig
	Http internalhttp.Config
	DB   db.Config
}

var config Config

func InitConfig() error {
	_ = godotenv.Overload()

	adder.SetConfigName("application")
	adder.AddConfigPath(".")
	adder.SetConfigType("yaml")
	adder.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	adder.AutomaticEnv()

	// Bind Secret Environment Variables
	//_ = adder.BindEnv("auth.secretkey", "OPENAI_SECRET_KEY")

	if err := adder.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := adder.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	initLogger(config.Log.Level)

	if strings.ToUpper(config.Log.Level) == LOG_LEVEL_DEBUG {
		configJSON, err := adder.PrettyJSON(config)
		if err == nil {
			slog.Debug("Config loaded:")
			slog.Debug(configJSON)
		}
	}

	return nil
}

func (c Config) Validate() error {
	return nil
}
