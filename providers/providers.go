package providers

import (
	"github.com/jmoiron/sqlx"
)

type RedisProvider interface {
	Get() (string, string, error)
	Publish(key string, value interface{}) error
}

type DBProvider interface {
	Ping() error
	Tx(fn func(tx *sqlx.Tx) error) error
	PSQLProvider
}

type PSQLProvider interface {
	DB() *sqlx.DB
}
