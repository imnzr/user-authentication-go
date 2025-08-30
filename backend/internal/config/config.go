package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type DatabaseConfig struct {
	// Primary database (MySQL)
	Primary MySQLConnection `json:"primary"`
	// Read replica
	ReadReplica *MySQLConnection `json:"read_replica,omitempty"`
	// Connection pool settings
	Pool PoolConfig `json:"pool"`
	// Mysql spesific settings
	MySQL MySQLConfig `json:"mysql"`
	// Migration settings
	Migration MigrationConfig `json:"migration"`
}

type MySQLConnection struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string

	// MySQL spesific connection parameters
	Charset   string
	Collation string
	ParseTime bool
	Loc       string
	TLS       string
}

type MySQLConfig struct {
	// Connection timeout settings
	Timeout      time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// MySQL spesific settings
	MultiStatements   bool
	InterpolateParams bool

	// Transaction isolation levels
	TxIsolation string
}

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type MigrationConfig struct {
	Enabled   bool
	Directory string
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Build MySQL DSN (data source name)
func (m *MySQLConnection) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&collation=%s&parseTime=%t&loc=%s&tls=%s",
		m.Username,
		m.Password,
		m.Password,
		m.Host,
		m.Database,
		m.Charset,
		m.Collation,
		m.ParseTime,
		m.Loc,
		m.TLS,
	)
}

// Build MySQL DSN with additional parameteres
func (m *MySQLConnection) DSNWithParams(mysqlConfig MySQLConfig) string {
	baseDSN := m.DSN()

	// Add MySQL spesific parameters
	params := fmt.Sprintf(
		"&timeout=%s&readTimeout=%s&writeTimeout=%s&multiStatements=%t&interpolateParams=%t",
		mysqlConfig.Timeout,
		mysqlConfig.ReadTimeout,
		mysqlConfig.WriteTimeout,
		mysqlConfig.MultiStatements,
		mysqlConfig.InterpolateParams,
	)

	return baseDSN + params
}

// Load MySQL configuration from environment
func loadDatabaseConfig(cfg *DatabaseConfig) error {
	// Primary database
	cfg.Primary = MySQLConnection{
		Host:      getEnvOrDefault("DB_HOST", "localhost"),
		Port:      os.Getenv("DB_PORT"),
		Database:  os.Getenv("DB_NAME"),
		Username:  os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASS"),
		Charset:   getEnvOrDefault("DB_CHARSET", "utf8mb4"),
		Collation: getEnvOrDefault("DB_COLLATION", "utf8mb4_unicode_ci"),
		ParseTime: getEnvBoolOrDefault("DB_PARSE_TIME", true),
		Loc:       getEnvOrDefault("DB_LOC", "UTC"),
		TLS:       getEnvOrDefault("DB_TLS", "false"),
	}
	// Validate required fields
	if cfg.Primary.Database == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if cfg.Primary.Username == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if cfg.Primary.Password == "" {
		return fmt.Errorf("DB_PASS is required")
	}

	// MySQL spesific configuration
	cfg.MySQL = MySQLConfig{
		Timeout:           getEnvDurationOrDefault("DB_TIMEOUT", 10*time.Second),
		ReadTimeout:       getEnvDurationOrDefault("DB_READ_TIMEOUT", 30*time.Second),
		WriteTimeout:      getEnvDurationOrDefault("DB_WRITE_TIMEOUT", 30*time.Second),
		MultiStatements:   getEnvBoolOrDefault("DB_MULTI_STATEMENTS", false),
		InterpolateParams: getEnvBoolOrDefault("DB_INTERPOLATE_PARAMS", true),
		TxIsolation:       getEnvOrDefault("DB_TX_ISOLATION", "READ-COMMITED"),
	}
	// Pool Configuration
	cfg.Pool = PoolConfig{
		MaxOpenConns:    getEnvIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvIntOrDefault("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Second),
		ConnMaxIdleTime: getEnvDurationOrDefault("DB_CONN_MAX_IDLE_TIME", 1*time.Second),
	}
	// Migration configuration
	cfg.Migration = MigrationConfig{
		Enabled:   getEnvBoolOrDefault("DB_MIGRATION_ENABLED", true),
		Directory: getEnvOrDefault("DB_MIGRATION_DIR", "migrations"),
	}

	return nil
}
