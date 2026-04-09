package database

import (
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"

	"billsoftware/backend/internal/config"
)

func NewMySQLEngine(cfg config.DatabaseConfig) (*xorm.Engine, error) {
	loc := url.QueryEscape(cfg.Loc)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
		cfg.ParseTime,
		loc,
	)

	engine, err := xorm.NewEngine(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	engine.SetMaxOpenConns(cfg.MaxOpenConns)
	engine.SetMaxIdleConns(cfg.MaxIdleConns)

	if err := engine.Ping(); err != nil {
		return nil, err
	}

	return engine, nil
}

