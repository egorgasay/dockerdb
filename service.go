// Package dockerdb allows user to create virtual databases using docker.
package dockerdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strings"
	"time"
)

const (
	tryInterval = 1 * time.Second

	//MSSQL      = "mssql"
	//MSSQLImage = "mcr.microsoft.com/mssql/server"

	Postgres      = "postgres"
	PostgresImage = "postgres"

	MySQL      = "mysql"
	MySQLImage = "mysql"
)

var (
	maxWaitTime          = 20 * time.Second
	ErrUnsupportedVendor = errors.New("following vendor is unsupported")
	ErrUnknown           = errors.New("unknown error")
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
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

inner:
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Trim(name, "/") == conf.DB.Name {
				ddb.ID = container.ID
				break inner
			}
		}
	}

	if ddb.ID != "" {
		err = ddb.setup(ctx)
		if err != nil {
			return nil, err
		}

		return ddb, nil
	}

	err = ddb.init(ctx)
	if err != nil {
		return ddb, err
	}

	return ddb, nil
}
