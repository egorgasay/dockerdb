package dockerdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (ddb *VDB) setup(ctx context.Context) error {
	err := ddb.Run(ctx)
	if err != nil {
		return err
	}

	if ddb.conf.noSQL {
		return nil
	}

	if ddb.connStr == "" {
		ddb.connStr, err = buildConnStr(ddb.conf)
		if err != nil {
			return err
		}
	}

	ddb.db, err = ddb.getDB(ddb.connStr)
	if err != nil {
		return err
	}

	return nil
}

func (ddb *VDB) getDB(connStr string) (db *sql.DB, err error) {
	after := time.After(_maxWaitTime)
	ticker := time.NewTicker(tryInterval)
	var errPing error
	for {
		select {
		case <-after:
			if errPing != nil {
				return nil, errPing
			}
			return nil, fmt.Errorf("timeout, Last error:%w", err)
		default:
			db, err = sql.Open(ddb.conf.vendorName, connStr)
			if db == nil {
				return nil, fmt.Errorf("db is nil %w", err)
			}

			errPing = db.Ping()
			if errPing == nil && err == nil {
				return db, nil
			}

			<-ticker.C
		}
	}
}
