package config

import (
	"fmt"
	"os"
	"time"
)

type DatabaseConfig struct {
	Primary   MySQLConnection `json:"primary"`
	MySQL     MySQLConfig     `json:"mysql"`
	Pool      PoolConfig      `json:"pool"`
	Migration MigrationConfig `json:"migration"`
}

type MySQLConnection struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Database  string `json:"database"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Charset   string `json:"charset"`
	Collation string `json:"collation"`
	ParseTime bool   `json:"parse_time"`
	Loc       string `json:"loc"`
	TLS       string `json:"tls"`
}

type MySQLConfig struct {
	Timeout           time.Duration `json:"timeout"`
	ReadTimeout       time.Duration `json:"read_timeout"`
	WriteTimeout      time.Duration `json:"write_timeout"`
	MultiStatements   bool          `json:"multi_statements"`
	InterpolateParams bool          `json:"interpolate_params"`
	TxIsolation       string        `json:"tx_isolation"`
}

type PoolConfig struct {
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

type MigrationConfig struct {
	Enabled   bool   `json:"enabled"`
	Directory string `json:"directory"`
}

// Build MySQL DSN (Data Source Name)
func (m *MySQLConnection) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=%t&loc=%s&tls=%s",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
		m.Charset,
		m.Collation,
		m.ParseTime,
		m.Loc,
		m.TLS,
	)
}

// Build MySQL DSN with additional parameters
func (m *MySQLConnection) DSNWithParams(mysqlCfg MySQLConfig) string {
	baseDSN := m.DSN()

	params := fmt.Sprintf(
		"&timeout=%s&readTimeout=%s&writeTimeout=%s&multiStatements=%t&interpolateParams=%t",
		mysqlCfg.Timeout,
		mysqlCfg.ReadTimeout,
		mysqlCfg.WriteTimeout,
		mysqlCfg.MultiStatements,
		mysqlCfg.InterpolateParams,
	)

	return baseDSN + params
}

// Load MySQL configuration from environment
func loadDatabaseConfig(cfg *DatabaseConfig) error {
	// Primary database connection
	cfg.Primary = MySQLConnection{
		Host:      getEnvOrDefault("DB_HOST", "localhost"),
		Port:      getEnvIntOrDefault("DB_PORT", 3306),
		Database:  os.Getenv("DB_NAME"),
		Username:  os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASSWORD"),
		Charset:   getEnvOrDefault("DB_CHARSET", "utf8mb4"),
		Collation: getEnvOrDefault("DB_COLLATION", "utf8mb4_unicode_ci"),
		ParseTime: getEnvBoolOrDefault("DB_PARSE_TIME", true),
		Loc:       getEnvOrDefault("DB_LOCATION", "UTC"),
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
		return fmt.Errorf("DB_PASSWORD is required")
	}

	// MySQL specific configuration
	cfg.MySQL = MySQLConfig{
		Timeout:           getEnvDurationOrDefault("DB_TIMEOUT", 10*time.Second),
		ReadTimeout:       getEnvDurationOrDefault("DB_READ_TIMEOUT", 30*time.Second),
		WriteTimeout:      getEnvDurationOrDefault("DB_WRITE_TIMEOUT", 30*time.Second),
		MultiStatements:   getEnvBoolOrDefault("DB_MULTI_STATEMENTS", false),
		InterpolateParams: getEnvBoolOrDefault("DB_INTERPOLATE_PARAMS", true),
		TxIsolation:       getEnvOrDefault("DB_TX_ISOLATION", "READ-COMMITTED"),
	}

	// Pool configuration
	cfg.Pool = PoolConfig{
		MaxOpenConns:    getEnvIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvIntOrDefault("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getEnvDurationOrDefault("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
	}

	// Migration configuration
	cfg.Migration = MigrationConfig{
		Enabled:   getEnvBoolOrDefault("DB_MIGRATION_ENABLED", true),
		Directory: getEnvOrDefault("DB_MIGRATION_DIR", "migrations"),
	}

	return nil
}
