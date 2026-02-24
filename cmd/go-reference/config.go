package main

import (
	"encoding/json"
	"fmt"
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

func InitConfig() {
	_ = godotenv.Load()

	adder.SetConfigName("application")
	adder.AddConfigPath(".")
	adder.SetConfigType("yaml")
	adder.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	adder.AutomaticEnv()

	if err := adder.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := adder.Unmarshal(&config); err != nil {
		panic(err)
	}

	initLogger(config.Log.Level)

	if strings.ToUpper(config.Log.Level) == LOG_LEVEL_DEBUG {
		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err == nil {
			fmt.Println("Config loaded:")
			fmt.Println(string(configJSON))
		}
	}
}
