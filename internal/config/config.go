package config

import (
	"os"

	"github.com/google/go-github/v57/github"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	TelegramBot struct {
		Token string `yaml:"token"`
	} `yaml:"telegram_bot"`
	RepoHostingApis struct {
		GithubToken string `yaml:"github_token"`
	} `yaml:"repo_hosting_apis"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"db"`
	} `yaml:"postgres"`
	ZapConfig zap.Config
}

type RepoHostingClients struct {
	GitHub *github.Client
}

func NewConfig(filename string) (*AppConfig, error) {
	var config AppConfig

	configFile, err := os.ReadFile(filename)
	if err != nil {
		return &config, err
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	config.ZapConfig = zapConfig

	err = yaml.Unmarshal(configFile, &config)
	return &config, err
}
