package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port            int
	Environment     string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.MaxOpenConns,
	)
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix("RH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("PORT", 8080)
	v.SetDefault("ENV", "development")
	v.SetDefault("SHUTDOWN_TIMEOUT", "30s")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "reliabilityhub")
	v.SetDefault("DB_PASSWORD", "reliabilityhub")
	v.SetDefault("DB_NAME", "reliabilityhub")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "5m")
	v.SetDefault("REDIS_ADDR", "localhost:6379")

	cfg := &Config{
		Server: ServerConfig{
			Port:            v.GetInt("PORT"),
			Environment:     v.GetString("ENV"),
			ShutdownTimeout: v.GetDuration("SHUTDOWN_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetInt("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			Name:            v.GetString("DB_NAME"),
			SSLMode:         v.GetString("DB_SSL_MODE"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
	}
	return cfg, nil
}

func (c *Config) IsProd() bool {
	return c.Server.Environment == "production"
}
