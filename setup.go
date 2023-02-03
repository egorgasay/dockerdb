package dockerdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (ddb *VDB) setup(ctx context.Context) error {
	err := ddb.Run(ctx)
	if err != nil {
		return err
	}

	ddb.ConnString = Build(ddb.conf)

	ddb.DB, err = ddb.getDB(ddb.ConnString)
	if err != nil {
		return err
	}

	return nil
}

func (ddb *VDB) getDB(connStr string) (*sql.DB, error) {
	after := time.After(maxWaitTime)
	ticker := time.NewTicker(tryInterval)
	for {
		select {
		case <-after:
			return nil, errors.New("timeout")
		default:
			db, err := sql.Open(ddb.conf.Vendor, connStr)
			if db == nil {
				return nil, fmt.Errorf("DB is nil %w", err)
			}

			pingErr := db.Ping()
			if pingErr == nil && err == nil {
				return db, nil
			}

			<-ticker.C
		}
	}
}
