package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	RocketMQ RocketMQConfig `yaml:"rocketmq"`
	JWT      JWTConfig      `yaml:"jwt"`
	CORS     CORSConfig     `yaml:"cors"`
	Media    MediaConfig    `yaml:"media"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	DSN                   string `yaml:"dsn"`
	MaxIdleConns          int    `yaml:"max_idle_conns"`
	MaxOpenConns          int    `yaml:"max_open_conns"`
	ConnMaxLifetimeMinute int    `yaml:"conn_max_lifetime_minute"`
	ConnectTimeoutSecond  int    `yaml:"connect_timeout_second"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type RocketMQConfig struct {
	Enabled             bool   `yaml:"enabled"`
	NameServerAddr      string `yaml:"name_server_addr"`
	Namespace           string `yaml:"namespace"`
	TopicPrefix         string `yaml:"topic_prefix"`
	ProducerGroup       string `yaml:"producer_group"`
	ConsumerGroupPrefix string `yaml:"consumer_group_prefix"`
	DialTimeoutSecond   int    `yaml:"dial_timeout_second"`
}

type JWTConfig struct {
	AccessSecret    string `yaml:"access_secret"`
	RefreshSecret   string `yaml:"refresh_secret"`
	Issuer          string `yaml:"issuer"`
	AccessTTLMinute int    `yaml:"access_ttl_minute"`
	RefreshTTLHour  int    `yaml:"refresh_ttl_hour"`
}

type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

type MediaConfig struct {
	RootDir       string `yaml:"root_dir"`
	PublicBaseURL string `yaml:"public_base_url"`
}

func Load(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file %q: %w", path, err)
	}

	cfg.applyDefaults()
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config file %q: %w", path, err)
	}
	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}
	if c.Database.MaxIdleConns == 0 {
		c.Database.MaxIdleConns = 5
	}
	if c.Database.MaxOpenConns == 0 {
		c.Database.MaxOpenConns = 10
	}
	if c.Database.ConnMaxLifetimeMinute == 0 {
		c.Database.ConnMaxLifetimeMinute = 30
	}
	if c.Database.ConnectTimeoutSecond == 0 {
		c.Database.ConnectTimeoutSecond = 5
	}
	if c.JWT.Issuer == "" {
		c.JWT.Issuer = "pilipili-go"
	}
	if c.JWT.AccessTTLMinute == 0 {
		c.JWT.AccessTTLMinute = 120
	}
	if c.JWT.RefreshTTLHour == 0 {
		c.JWT.RefreshTTLHour = 24 * 14
	}
	if len(c.CORS.AllowOrigins) == 0 {
		c.CORS.AllowOrigins = defaultCORSOrigins()
		c.CORS.AllowCredentials = true
	}
	if strings.TrimSpace(c.RocketMQ.NameServerAddr) == "" {
		c.RocketMQ.NameServerAddr = "127.0.0.1:9876"
	}
	if strings.TrimSpace(c.RocketMQ.TopicPrefix) == "" {
		c.RocketMQ.TopicPrefix = "pilipili-go"
	}
	if strings.TrimSpace(c.RocketMQ.ProducerGroup) == "" {
		c.RocketMQ.ProducerGroup = "pilipili-go-api"
	}
	if strings.TrimSpace(c.RocketMQ.ConsumerGroupPrefix) == "" {
		c.RocketMQ.ConsumerGroupPrefix = "pilipili-go-worker"
	}
	if c.RocketMQ.DialTimeoutSecond == 0 {
		c.RocketMQ.DialTimeoutSecond = 3
	}
	if strings.TrimSpace(c.Media.RootDir) == "" {
		c.Media.RootDir = "storage"
	}
	if strings.TrimSpace(c.Media.PublicBaseURL) == "" {
		c.Media.PublicBaseURL = "/uploads"
	}
}

func (c *Config) validate() error {
	if strings.TrimSpace(c.Database.DSN) == "" {
		return fmt.Errorf("database.dsn is required")
	}
	if strings.TrimSpace(c.Redis.Addr) == "" {
		return fmt.Errorf("redis.addr is required")
	}
	if c.RocketMQ.Enabled {
		if strings.TrimSpace(c.RocketMQ.NameServerAddr) == "" {
			return fmt.Errorf("rocketmq.name_server_addr is required when rocketmq.enabled=true")
		}
		if c.RocketMQ.DialTimeoutSecond <= 0 {
			return fmt.Errorf("rocketmq.dial_timeout_second must be greater than 0")
		}
	}
	if err := validateSecret("jwt.access_secret", c.JWT.AccessSecret); err != nil {
		return err
	}
	if err := validateSecret("jwt.refresh_secret", c.JWT.RefreshSecret); err != nil {
		return err
	}
	if c.Server.Port <= 0 {
		return fmt.Errorf("server.port must be greater than 0")
	}
	if c.JWT.AccessTTLMinute <= 0 {
		return fmt.Errorf("jwt.access_ttl_minute must be greater than 0")
	}
	if c.JWT.RefreshTTLHour <= 0 {
		return fmt.Errorf("jwt.refresh_ttl_hour must be greater than 0")
	}
	if strings.TrimSpace(c.Media.RootDir) == "" {
		return fmt.Errorf("media.root_dir is required")
	}
	if !strings.HasPrefix(strings.TrimSpace(c.Media.PublicBaseURL), "/") {
		return fmt.Errorf("media.public_base_url must start with /")
	}

	return nil
}

func validateSecret(field string, value string) error {
	secret := strings.TrimSpace(value)
	switch {
	case secret == "":
		return fmt.Errorf("%s is required", field)
	case strings.Contains(secret, "change-me"):
		return fmt.Errorf("%s must not use placeholder values", field)
	default:
		return nil
	}
}

func defaultCORSOrigins() []string {
	return []string{
		"http://127.0.0.1:3000",
		"http://localhost:3000",
		"http://127.0.0.1:5173",
		"http://localhost:5173",
		"http://127.0.0.1:8081",
		"http://localhost:8081",
	}
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
