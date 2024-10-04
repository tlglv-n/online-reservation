package store

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"
)

type SQLX struct {
	Client *sqlx.DB
}

func NewSQL(dataSourceName string) (store SQLX, err error) {
	if !strings.Contains(dataSourceName, "://") {
		err = errors.New("store: undefined data source name " + dataSourceName)
		return
	}
	driverName := strings.ToLower(strings.Split(dataSourceName, "://")[0])

	store.Client, err = sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return
	}
	store.Client.SetMaxOpenConns(20)

	return
}
