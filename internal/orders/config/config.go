package config

import (
	"github.com/spf13/viper"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

type (
	Config struct {
		Env      string   `mapstructure:"environment"` // лучше быть в переменной окружения
		HTTP     HTTP     `mapstructure:"http"`
		Postgres Postgres `mapstructure:"postgres"`
		Kafka    Kafka    `mapstructure:"kafka"`
	}

	HTTP struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	Kafka struct {
		Topic              string   `mapstructure:"topic"`
		ConsumerGroup      string   `mapstructure:"consumer_group"`
		SessionTimeout     int      `mapstructure:"session_timeout"`
		AutoCommitInterval int      `mapstructure:"auto_commit_inteval"`
		BootstrapServers   []string `mapstructure:"bootstrap_servers"`
	}
)

func MustLoad(folder string) (*Config, error) {

	if err := parseConfigFile(folder); err != nil {
		return nil, err
	}

	var config Config

	if err := unmarshall(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func unmarshall(cfg *Config) error {
	if err := viper.UnmarshalKey("environment", &cfg.Env); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("kafka", &cfg.Kafka); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	env := viper.GetString("environment")
	if env == envLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()

}
