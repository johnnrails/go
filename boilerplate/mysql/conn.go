package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/johnnrails/ddd_go/boilerplate/logger"
)

type ConnectionStruct struct {
	Host            string
	Port            int
	User            string
	Pass            string
	Database        string
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

func NewConnection(ctx context.Context, cfg ConnectionStruct) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Database))
	if err != nil {
		logger.Critical(ctx, fmt.Sprintf("MYSQL CONNECTION ERROR: %s", err))
	}
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	return db
}
