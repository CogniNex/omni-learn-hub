package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App    `yaml:"app"`
		HTTP   `yaml:"http"`
		Log    `yaml:"logger"`
		PG     `yaml:"postgres"`
		OTP    `yaml:"otp"`
		SMS    `yaml:"sms-service"`
		Vonage `yaml:"vonage-client"`
		JWT    `yaml:"JWT"`
	}
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `                                    env:"PG_URL"`
	}

	OTP struct {
		Length int `env-required:"true" yaml:"length" env:"LENGTH"`
	}

	SMS struct {
		Templates Templates `env-required:"true" yaml:"templates"`
		From      string    `env-requred:"true" yaml:"from"`
	}

	Vonage struct {
		ApiKey    string `env-required:"true" yaml:"api_key"`
		ApiSecret string `env-required:"true" yaml:"api_secret"`
	}

	Templates struct {
		Registration string `env-required:"true" yaml:"registration"`
	}

	JWT struct {
		AccessTokenExpireDuration  time.Duration `env-required:"true" yaml:"accessTokenExpireDuration"`
		RefreshTokenExpireDuration time.Duration `env-required:"true" yaml:"refreshTokenExpireDuration"`
		Secret                     string        `env-required:"true" yaml:"secret"`
		RefreshSecret              string        `env-required:"true" yaml:"refreshSecret"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	cfg, err := ParseConfigFiles("./config/config-prod.yml", ".env")
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ParseConfigFiles(files ...string) (*Config, error) {
	var cfg Config

	for i := 0; i < len(files); i++ {
		err := cleanenv.ReadConfig(files[i], &cfg)
		if err != nil {
			log.Printf("Error reading configuration from file:%v", files[i])
			return nil, err
		}
	}

	return &cfg, nil
}
