package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imnzr/user-authentication-go/internal/config"
)

type DB struct {
	Primary *sql.DB
	Config  *config.DatabaseConfig
}

// New initializes MySQL database connection
func New(cfg *config.DatabaseConfig) (*DB, error) {
	// Connect to primary database
	primary, err := connectMySQL(cfg.Primary, cfg.MySQL, cfg.Pool)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	db := &DB{
		Primary: primary,
		Config:  cfg,
	}

	// Configure MySQL session settings
	if err := db.configureMySQLSession(); err != nil {
		return nil, fmt.Errorf("failed to configure MySQL session: %w", err)
	}

	return db, nil
}

func connectMySQL(conn config.MySQLConnection, mysqlCfg config.MySQLConfig, pool config.PoolConfig) (*sql.DB, error) {
	dsn := conn.DSNWithParams(mysqlCfg)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(pool.MaxOpenConns)
	db.SetMaxIdleConns(pool.MaxIdleConns)
	db.SetConnMaxLifetime(pool.ConnMaxLifetime)
	db.SetConnMaxIdleTime(pool.ConnMaxIdleTime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MySQL database: %w", err)
	}

	return db, nil
}

func (db *DB) configureMySQLSession() error {
	// Set transaction isolation level
	query := fmt.Sprintf("SET SESSION transaction_isolation = '%s'", db.Config.MySQL.TxIsolation)
	if _, err := db.Primary.Exec(query); err != nil {
		return fmt.Errorf("failed to set transaction isolation: %w", err)
	}

	// Set SQL mode for strict error handling
	// sqlMode := "STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION"
	// if _, err := db.Primary.Exec("SET SESSION sql_mode = ?", sqlMode); err != nil {
	// 	return fmt.Errorf("failed to set SQL mode: %w", err)
	// }

	// Set timezone to UTC
	if _, err := db.Primary.Exec("SET time_zone = '+00:00'"); err != nil {
		return fmt.Errorf("failed to set timezone: %w", err)
	}

	return nil
}

// Close closes database connection
func (db *DB) Close() error {
	if db.Primary != nil {
		return db.Primary.Close()
	}
	return nil
}

// Health checks database connectivity
func (db *DB) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.Primary.PingContext(ctx)
}

// GetVersion returns MySQL version
func (db *DB) GetVersion() (string, error) {
	var version string
	err := db.Primary.QueryRow("SELECT VERSION()").Scan(&version)
	return version, err
}

// GetStats returns database statistics
func (db *DB) GetStats() sql.DBStats {
	return db.Primary.Stats()
}
