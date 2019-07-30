package storage

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/partyzanex/esender/domain"
	"github.com/partyzanex/esender/storage/mysql"
	"github.com/pkg/errors"
)

type Config struct {
	Name         string
	DSN          string
	ConnLifetime time.Duration
}

func Create(config Config) (domain.EmailStorage, error) {
	switch config.Name {
	case "mysql":
		db, err := sql.Open("mysql", config.DSN)
		if err != nil {
			return nil, errors.Wrap(err, "opening sql connection failed")
		}
		if config.ConnLifetime > 0 {
			db.SetConnMaxLifetime(config.ConnLifetime)
		}

		return mysql.EmailStorage(db), nil
	default:
		return nil, errors.New("unknown storage name " + config.Name)
	}
}
