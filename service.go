// Package dockerdb allows user to create virtual databases using docker.
// Tested with PostgreSQL, MySQL, MS SQL.
package dockerdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"strconv"
	"strings"
	"time"
)

type DockerHubName string

const (
	tryInterval = 1 * time.Second

	Postgres15 DockerHubName = "postgres:15"
	Postgres14 DockerHubName = "postgres:14"
	Postgres13 DockerHubName = "postgres:13"
	Postgres12 DockerHubName = "postgres:12"
	Postgres11 DockerHubName = "postgres:11"

	MySQL5Image DockerHubName = "mysql:5.7"
	MySQL8Image DockerHubName = "mysql:8"

	KeyDBImage DockerHubName = "eqalpha/keydb"

	RedisImage DockerHubName = "redis"

	ScyllaDBImage DockerHubName = "scylladb/scylla"
)

const (
	postgres = "postgres"
	mysql    = "mysql"
)

var (
	maxWaitTime = 20 * time.Second
	ErrUnknown  = errors.New("unknown error")
)

type VDB struct {
	id      string
	cli     *client.Client
	conf    Config
	db      *sql.DB
	connStr string
}

// New creates a new docker container and launches it
func New(ctx context.Context, conf Config) (vdb *VDB, err error) {
	if conf.pullImage {
		ctx := context.TODO()
		err := Pull(ctx, conf.vendor)
		if err != nil {
			return nil, fmt.Errorf("pull error: %w", err)
		}
	}

	if err = validate(conf); err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	vdb = &VDB{}
	vdb.cli = cli
	vdb.conf = conf

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	vendor := strings.Split(string(vdb.conf.vendor), ":")
	if len(vendor) == 0 {
		return nil, errors.New("vendor must be not empty")
	}
	vdb.conf.vendorName = vendor[0]

inner:
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Trim(name, "/") == conf.db.Name {
				vdb.id = container.ID
				vdb.conf.actualPort = nat.Port(strconv.Itoa(int(container.Ports[0].PublicPort)))
				break inner
			}
		}
	}

	if vdb.id != "" {
		err = vdb.setup(ctx)
		if err != nil {
			return vdb, err
		}
	} else {
		if vdb.conf.actualPort == "" {
			vdb.conf.actualPort, err = getFreePort()
			if err != nil {
				return nil, err
			}
		}
		err = vdb.init(ctx)
		if err != nil {
			return vdb, err
		}
	}

	if vdb.conf.noSQL {
		check := vdb.conf.checkWakeUp
		var stop bool
		for i := 0; i < check.tries; i++ {
			stop = check.fn(vdb.conf)
			if stop {
				break
			}
			time.Sleep(check.sleepTime)
		}

		if !stop {
			return vdb, ErrUnknown
		}

		return vdb, nil
	}

	return vdb, nil
}
