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

type CustomDB struct {
	DB     DB
	Port   string
	Vendor string
}

func New(ctx context.Context, conf CustomDB) (*VDB, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	ddb := &VDB{cli: cli, conf: conf}
	err = ddb.pull(ctx, ddb.conf.Vendor)
	if err != nil {
		return nil, err
	}

	err = ddb.init(ctx)
	if err != nil {
		return nil, err
	}

	return ddb, nil
}
