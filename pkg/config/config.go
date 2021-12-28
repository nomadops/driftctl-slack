package viper

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config is a struct that contains the configuration environment for the application.
type Config struct {
	ScanReport  string
	ScanBucket  string
	StateBucket string
	Channel     string
	Token       string
}

// LoadConfig loads the configuration from the environment.
func LoadConfig() (env *Config, err error) {
	var config Config
	v := viper.New()
	v.SetConfigType("env")
	v.BindEnv("scanReport", "SCAN_FILE")
	v.BindEnv("scanBucket", "SCAN_BUCKET")
	v.BindEnv("stateBucket", "STATE_BUCKET")
	v.BindEnv("channel", "CHANNEL")
	v.BindEnv("token", "TOKEN")

	for k, v := range v.AllSettings() {
		fmt.Println(v)
		if v == "" {
			log.Fatal().
				Str("service", "driftctl-slack").
				Msg("Config value for " + k + " is not set.")
		}
	}

	err = v.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}
	return &config, nil
}
