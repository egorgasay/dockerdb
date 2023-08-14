package dockerdb

import (
	"context"
	"database/sql"
	"github.com/docker/go-connections/nat"
)

// SQL returns an *sql.DB instance.
func (ddb *VDB) SQL() (db *sql.DB) {
	return ddb.db
}

func (ddb *VDB) GetPort() (port nat.Port) {
	return ddb.conf.actualPort
}

func With(ctx context.Context, c Config, fn func(c Config, vdb *VDB) error) (err error) {
	vdb, err := New(ctx, c)
	if err != nil {
		if vdb != nil {
			return vdb.Clear(ctx)
		}
		return err
	}
	defer func() { err = vdb.Clear(ctx) }()

	return fn(c, vdb)
}

func WithPostgres(ctx context.Context, name string, fn func(c Config, vdb *VDB) error) (err error) {
	return With(ctx, PostgresConfig(name).Build(), fn)
}

func WithMySQL(ctx context.Context, name string, fn func(c Config, vdb *VDB) error) (err error) {
	return With(ctx, MySQLConfig(name).Build(), fn)
}
