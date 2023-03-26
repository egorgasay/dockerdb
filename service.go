// Package dockerdb allows user to create virtual databases using docker.
// Tested with PostgreSQL, MySQL, MS SQL.
package dockerdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"strings"
	"time"
)

const (
	tryInterval = 1 * time.Second

	Postgres15 = "postgres:15"
	Postgres14 = "postgres:14"
	Postgres13 = "postgres:13"
	Postgres12 = "postgres:12"
	Postgres11 = "postgres:11"

	MySQL5Image = "mysql:5.7"
	MySQL8Image = "mysql:8"
)

var (
	maxWaitTime    = 20 * time.Second
	ErrUnknown     = errors.New("unknown error")
	ErrUnsupported = errors.New("unsupported db vendor")
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

type CustomDB struct {
	DB         DB
	Port       string
	Vendor     string
	vendorName string

	// Optional if you are using a supported vendor
	PortDocker nat.Port
	EnvDocker  []string
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

	vendor := strings.Split(ddb.conf.Vendor, ":")
	if len(vendor) == 0 {
		return nil, errors.New("vendor must be not empty")
	}
	ddb.conf.vendorName = vendor[0]

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
