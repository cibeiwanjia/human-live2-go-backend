// Package config provides configuration management using Viper
package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Common   CommonConfig   `mapstructure:"common"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Agents   AgentsConfig   `mapstructure:"agents"`
	Engines  EnginesConfig  `mapstructure:"engines"`
}

// CommonConfig holds common application settings
type CommonConfig struct {
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	LogLevel string `mapstructure:"log_level"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	IP            string `mapstructure:"ip"`
	Port          int    `mapstructure:"port"`
	WorkspacePath string `mapstructure:"workspace_path"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// AgentsConfig holds agent configuration
type AgentsConfig struct {
	SupportList []AgentItemConfig `mapstructure:"support_list"`
	Default     string            `mapstructure:"default"`
}

// AgentItemConfig holds individual agent configuration
type AgentItemConfig struct {
	Name string `mapstructure:"name"`
	Type string `mapstructure:"type"`
	Desc string `mapstructure:"desc"`
}

// EnginesConfig holds engine configurations
type EnginesConfig struct {
	TTS TTSEnginesConfig `mapstructure:"tts"`
	ASR ASREnginesConfig `mapstructure:"asr"`
}

// TTSEnginesConfig holds TTS engine configuration
type TTSEnginesConfig struct {
	SupportList []EngineItemConfig `mapstructure:"support_list"`
	Default     string             `mapstructure:"default"`
}

// ASREnginesConfig holds ASR engine configuration
type ASREnginesConfig struct {
	SupportList []EngineItemConfig `mapstructure:"support_list"`
	Default     string             `mapstructure:"default"`
}

// EngineItemConfig holds individual engine configuration
type EngineItemConfig struct {
	Name string `mapstructure:"name"`
	Type string `mapstructure:"type"`
	Desc string `mapstructure:"desc"`
}

var (
	configOnce        sync.Once
	appConfig         *Config
	defaultConfigPath = "./configs/config.yaml"
)

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	var err error
	configOnce.Do(func() {
		appConfig, err = loadConfig(configPath)
	})
	return appConfig, err
}

// Get returns the loaded configuration
func Get() *Config {
	if appConfig == nil {
		// Load with default path if not loaded
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = defaultConfigPath
		}
		var err error
		appConfig, err = loadConfig(configPath)
		if err != nil {
			panic(fmt.Sprintf("failed to load config: %v", err))
		}
	}
	return appConfig
}

func loadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file
	v.SetConfigFile(configPath)

	// Set default values
	setDefaults(v)

	// Read environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables in config values
	expandEnvVars(v)

	// Unmarshal to struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with environment variables
	overrideWithEnv(&cfg)

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("common.name", "Awesome-Digital-Human")
	v.SetDefault("common.version", "v3.0.0")
	v.SetDefault("common.log_level", "info")
	v.SetDefault("server.ip", "0.0.0.0")
	v.SetDefault("server.port", 8881)
	v.SetDefault("server.workspace_path", "./outputs")
	v.SetDefault("database.type", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
}

// expandEnvVars expands environment variables in string values
// Supports ${VAR} and ${VAR:default} syntax
func expandEnvVars(v *viper.Viper) {
	for _, key := range v.AllKeys() {
		val := v.GetString(key)
		if strings.Contains(val, "${") {
			expanded := os.Expand(val, func(envKey string) string {
				// Handle ${VAR:default} syntax
				parts := strings.SplitN(envKey, ":", 2)
				envVal := os.Getenv(parts[0])
				if envVal == "" && len(parts) > 1 {
					return parts[1]
				}
				return envVal
			})
			v.Set(key, expanded)
		}
	}
}

func overrideWithEnv(cfg *Config) {
	// Override with explicit environment variables
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Common.LogLevel = logLevel
	}
	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		var port int
		fmt.Sscanf(serverPort, "%d", &port)
		if port > 0 {
			cfg.Server.Port = port
		}
	}
}
