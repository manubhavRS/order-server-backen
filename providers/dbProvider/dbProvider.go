package dbProvider

import (
	"OrderServer/models"
	"OrderServer/providers"
	"OrderServer/utils"
	"fmt"
	"log"

	// source/file import is required for migration files to read
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	// load pq as database driver
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type psqlProvider struct {
	db *sqlx.DB
}

func NewPSQLProvider(configs models.DatabaseConfig, sslMode utils.SSLMode) providers.DBProvider {
	//"host=localhost port=5432 user=local password=local dbname=audioPhile sslmode=disable"
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", configs.Host, configs.Port, configs.User, configs.Password, configs.DatabaseName, sslMode)
	log.Printf(connStr)
	DB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Panicf("error: %v", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Panicf("error: %v", err)
	}
	return &psqlProvider{
		db: DB,
	}
}
func (pp *psqlProvider) Ping() error {
	return pp.db.Ping()
}
func (pp *psqlProvider) DB() *sqlx.DB {
	return pp.db
}

// Tx provides the transaction wrapper
func (pp *psqlProvider) Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := pp.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}
