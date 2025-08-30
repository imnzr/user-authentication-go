package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imnzr/user-authentication-go/internal/config"
)

type DB struct {
	Primary     *sql.DB
	ReadReplica *sql.DB
	Config      *config.DatabaseConfig
}
