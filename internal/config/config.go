package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds application configuration loaded from configs/app.yaml and environment overrides.
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	HTTPClient HTTPClientConfig `mapstructure:"httpClient"`
	Egov       EgovConfig       `mapstructure:"egov"`
	Search     SearchConfig     `mapstructure:"search"`
}

type HTTPClientConfig struct {
	TimeoutSeconds int `mapstructure:"timeoutSeconds"`
}

type ServerConfig struct {
	Port        string `mapstructure:"port"`
	ContextPath string `mapstructure:"contextPath"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslMode"`
}

type KafkaConfig struct {
	BootstrapServers    []string              `mapstructure:"bootstrapServers"`
	SaveLandInfoTopic   string                `mapstructure:"saveLandInfoTopic"`
	UpdateLandInfoTopic string                `mapstructure:"updateLandInfoTopic"`
	Producer            KafkaProducerConfig   `mapstructure:"producer"`
}

type KafkaProducerConfig struct {
	Retries      int    `mapstructure:"retries"`
	RequiredAcks string `mapstructure:"requiredAcks"`
}

type EgovConfig struct {
	User              UserConfig              `mapstructure:"user"`
	Location          LocationConfig          `mapstructure:"location"`
	MDMS              MDMSConfig              `mapstructure:"mdms"`
	Localization      LocalizationConfig      `mapstructure:"localization"`
	Pagination        PaginationConfig        `mapstructure:"pagination"`
	OwnershipCategory OwnershipCategoryConfig `mapstructure:"ownershipCategory"`
}

type UserConfig struct {
	Host           string `mapstructure:"host"`
	ContextPath    string `mapstructure:"contextPath"`
	CreatePath     string `mapstructure:"createPath"`
	SearchPath     string `mapstructure:"searchPath"`
	UpdatePath     string `mapstructure:"updatePath"`
	UsernamePrefix string `mapstructure:"usernamePrefix"`
}

func (u UserConfig) CreateURL() string {
	return u.Host + u.ContextPath + u.CreatePath
}

func (u UserConfig) SearchURL() string {
	// Java LandUserService uses userHost + searchEndpoint only (no context path).
	return u.Host + u.SearchPath
}

func (u UserConfig) UpdateURL() string {
	return u.Host + u.ContextPath + u.UpdatePath
}

type LocationConfig struct {
	Host                string `mapstructure:"host"`
	ContextPath         string `mapstructure:"contextPath"`
	Endpoint            string `mapstructure:"endpoint"`
	HierarchyTypeCode   string `mapstructure:"hierarchyTypeCode"`
}

func (l LocationConfig) BoundarySearchURL() string {
	return l.Host + l.ContextPath + l.Endpoint
}

type MDMSConfig struct {
	Host           string `mapstructure:"host"`
	SearchEndpoint string `mapstructure:"searchEndpoint"`
}

func (m MDMSConfig) SearchURL() string {
	return m.Host + m.SearchEndpoint
}

type LocalizationConfig struct {
	Host           string `mapstructure:"host"`
	WorkDirPath    string `mapstructure:"workDirPath"`
	ContextPath    string `mapstructure:"contextPath"`
	SearchEndpoint string `mapstructure:"searchEndpoint"`
	StateLevel     bool   `mapstructure:"stateLevel"`
}

type PaginationConfig struct {
	DefaultOffset int `mapstructure:"defaultOffset"`
	DefaultLimit  int `mapstructure:"defaultLimit"`
	MaxLimit      int `mapstructure:"maxLimit"`
}

type OwnershipCategoryConfig struct {
	Institutional string `mapstructure:"institutional"`
}

type SearchConfig struct {
	CitizenAllowedParams  string `mapstructure:"citizenAllowedParams"`
	EmployeeAllowedParams string `mapstructure:"employeeAllowedParams"`
}

// DSN returns a PostgreSQL connection string for read-only search queries (Phase 4+).
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
	)
}

// Load reads configuration from configs/app.yaml (or path in LAND_CONFIG_FILE).
// Environment variables use prefix LAND_ with dots replaced by underscores, e.g. LAND_SERVER_PORT.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	v.AddConfigPath(".")
	v.SetEnvPrefix("LAND")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if path := os.Getenv("LAND_CONFIG_FILE"); path != "" {
		v.SetConfigFile(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.Server.Port == "" {
		cfg.Server.Port = "8199"
	}
	if cfg.Server.ContextPath == "" {
		cfg.Server.ContextPath = "/land-services"
	}

	return &cfg, nil
}
