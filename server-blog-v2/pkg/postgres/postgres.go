// Package postgres GORM Postgres 连接封装。
package postgres

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	_defaultMaxIdleConns = 10
	_defaultMaxOpenConns = 100
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres GORM 连接容器。
type Postgres struct {
	maxIdleConns int
	maxOpenConns int
	connAttempts int
	connTimeout  time.Duration

	DB *gorm.DB
}

// New 创建 Postgres 连接。
func New(host string, port int, user, password, dbname, sslmode, timeZone string, maxIdleConns, maxOpenConns int, opts ...Option) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", host, port, user, password, dbname, sslmode, timeZone)

	m := &Postgres{
		maxIdleConns: maxIdleConns,
		maxOpenConns: maxOpenConns,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	if m.maxIdleConns == 0 {
		m.maxIdleConns = _defaultMaxIdleConns
	}
	if m.maxOpenConns == 0 {
		m.maxOpenConns = _defaultMaxOpenConns
	}

	for _, opt := range opts {
		opt(m)
	}

	var err error
	for m.connAttempts > 0 {
		m.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := m.DB.DB()
			sqlDB.SetMaxOpenConns(m.maxOpenConns)
			sqlDB.SetMaxIdleConns(m.maxIdleConns)

			break
		}

		log.Printf("Postgres: connect retry, attempts left: %d", m.connAttempts)
		time.Sleep(m.connTimeout)
		m.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("Postgres - New - connAttempts == 0: %w", err)
	}

	return m, nil
}

// Close 关闭 Postgres 连接。
func (m *Postgres) Close() {
	if m.DB != nil {
		sqlDB, _ := m.DB.DB()
		sqlDB.Close()
		fmt.Println("Postgres: connection closed")
	}
}
