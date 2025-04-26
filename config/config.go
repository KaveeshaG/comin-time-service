package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	QR       QRConfig       `mapstructure:"qr"`
}

type ServerConfig struct {
	Port     string        `mapstructure:"port"`
	Timeouts TimeoutConfig `mapstructure:"timeouts"`
}

type TimeoutConfig struct {
	Read  int `mapstructure:"read"`
	Write int `mapstructure:"write"`
	Idle  int `mapstructure:"idle"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type QRConfig struct {
	Size int `mapstructure:"size"`
}

func LoadConfig(path string) (*Config, error) {
	// Set defaults
	viper.SetDefault("server.port", "8084")
	viper.SetDefault("server.timeouts.read", 10)
	viper.SetDefault("server.timeouts.write", 10)
	viper.SetDefault("server.timeouts.idle", 60)
	viper.SetDefault("database.max_open_conns", 20)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("qr.size", 256)

	// Set config file properties
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, using defaults and environment variables
		// Define database configuration if file not found
		viper.SetDefault("database.host", "ep-flat-shadow-a8onelva.eastus2.azure.neon.tech")
		viper.SetDefault("database.port", "5432")
		viper.SetDefault("database.user", "comin_owner")
		viper.SetDefault("database.password", "Ye5rfjcIB7FX")
		viper.SetDefault("database.dbname", "comin")
		viper.SetDefault("database.sslmode", "require")
		viper.SetDefault("database.max_open_conns", 20)
		viper.SetDefault("database.max_idle_conns", 5)
	}

	// Environment variables override configuration file
	// Database configuration from environment variables
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		viper.Set("database.host", dbHost)
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		viper.Set("database.port", dbPort)
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		viper.Set("database.user", dbUser)
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		viper.Set("database.password", dbPassword)
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		viper.Set("database.dbname", dbName)
	}
	if dbSSLMode := os.Getenv("DB_SSL_MODE"); dbSSLMode != "" {
		viper.Set("database.sslmode", dbSSLMode)
	}
	if maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConnsStr != "" {
		if maxOpenConns, err := strconv.Atoi(maxOpenConnsStr); err == nil {
			viper.Set("database.max_open_conns", maxOpenConns)
		}
	}
	if maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleConnsStr != "" {
		if maxIdleConns, err := strconv.Atoi(maxIdleConnsStr); err == nil {
			viper.Set("database.max_idle_conns", maxIdleConns)
		}
	}

	// Server configuration from environment variables
	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		viper.Set("server.port", serverPort)
	}

	// Parse the config into the struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return &config, nil
}
