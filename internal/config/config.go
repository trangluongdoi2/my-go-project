package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string `mapstructure:"host" yaml:"host"`
		Port     string `mapstructure:"port" yaml:"port"`
		Username string `mapstructure:"username" yaml:"username"`
		Password string `mapstructure:"password" yaml:"password"`
		DbName   string `mapstructure:"db_name" yaml:"db_name"`
	} `mapstructure:"database" yaml:"database"`
	Server struct {
		Port     string `mapstructure:"port" yaml:"port"`
		GRPCPort string `mapstructure:"grpc_port" yaml:"grpc_port"`
	} `mapstructure:"server" yaml:"server"`
	RabbitMQConfig struct {
		Schema    string            `mapstructure:"schema" yaml:"schema"`
		Worker    int               `mapstructure:"worker" yaml:"worker"`
		Host      string            `mapstructure:"host" yaml:"host"`
		Port      int               `mapstructure:"port" yaml:"port"`
		VHost     string            `mapstructure:"vhost" yaml:"vhost"`
		Username  string            `mapstructure:"username" yaml:"username"`
		Password  string            `mapstructure:"password" yaml:"password"`
		SSL       bool              `mapstructure:"ssl" yaml:"ssl"`
		CTag      string            `mapstructure:"ctag" yaml:"ctag"`
		Retry     int               `mapstructure:"retry" yaml:"retry"`
		Exchanges string            `mapstructure:"exchanges" yaml:"exchanges"`
		Queue     map[string]string `mapstructure:"queue" yaml:"queue"`
		Uri       string            `mapstructure:"uri" yaml:"uri"`
	} `mapstructure:"rabbitmq" yaml:"rabbitmq"`
	RedisConfig struct {
		Host     string `mapstructure:"host" yaml:"host"`
		Port     int    `mapstructure:"port" yaml:"port"`
		Password string `mapstructure:"password" yaml:"password"`
		DB       int    `mapstructure:"db" yaml:"db"`
	} `mapstructure:"redis" yaml:"redis"`
	JWTConfig struct {
		Secret                    string `mapstructure:"secret" yaml:"secret"`
		AccessTokenExpirationMin  int    `mapstructure:"access_token_expiration_min" yaml:"access_token_expiration_min"`
		RefreshTokenExpirationDay int    `mapstructure:"refresh_token_expiration_day" yaml:"refresh_token_expiration_day"`
		ExpirationHours           int    `mapstructure:"expiration_hours" yaml:"expiration_hours"`
	} `mapstructure:"jwt" yaml:"jwt"`
}

var config Config

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "localhost"
	}

	configPath, err := filepath.Abs("./internal/config")
	if err != nil {
		return nil, fmt.Errorf("error getting config path: %v", err)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(env)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return &config, nil
}

func GetConfig() *Config {
	return &config
}
