package storage

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/partyzanex/esender/domain"
	"github.com/pkg/errors"
	"github.com/partyzanex/esender/storage/mysql"
)

type Config struct {
	Name string
	DSN  string
}

func Create(config Config) (domain.EmailStorage, error) {
	switch config.Name {
	case "mysql":
		db, err := sql.Open("mysql", config.DSN)
		if err != nil {
			return nil, errors.Wrap(err, "opening sql connection failed")
		}

		return mysql.EmailStorage(db), nil
	default:
		return nil, errors.New("unknown storage name " + config.Name)
	}
}
