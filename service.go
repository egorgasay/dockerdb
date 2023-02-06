// Package dockerdb allows user to create virtual databases using docker.
package dockerdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/docker/docker/client"
	"time"
)

const (
	maxWaitTime = 20 * time.Second
	tryInterval = 1 * time.Second

	//MSSQL      = "mssql"
	//MSSQLImage = "mcr.microsoft.com/mssql/server"

	Postgres      = "postgres"
	PostgresImage = "postgres"

	MySQL      = "mysql"
	MySQLImage = "mysql"
)

var (
	ErrUnsupportedVendor = errors.New("following vendor is unsupported")
	ErrAlreadyBindPort   = errors.New("the port is already in use by another container")
)

type VDB struct {
	ID         string
	cli        *client.Client
	conf       CustomDB
	DB         *sql.DB
	ConnString string
}

type DB struct {
	Name     string
	User     string
	Password string
}

type Vendor struct {
	Name  string
	Image string
}

type CustomDB struct {
	DB     DB
	Port   string
	Vendor Vendor
}

// New creates a new docker container and launches it
func New(ctx context.Context, conf CustomDB) (*VDB, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	ddb := &VDB{cli: cli, conf: conf}

	err = ddb.init(ctx)
	if err != nil {
		return ddb, err
	}

	return ddb, nil
}
