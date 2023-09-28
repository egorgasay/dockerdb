package dockerdb

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"strings"
)

// init initializes the docker container with the selected database
func (ddb *VDB) init(ctx context.Context) error {
	var env []string
	var portDocker nat.Port
	var specifiedVendor = string(ddb.conf.vendor)

	if ddb.conf.actualPort == "" {
		return errors.New("port must be not empty")
	}

	vendor := strings.Split(specifiedVendor, ":")
	if len(vendor) == 0 {
		return errors.New("vendor must be not empty")
	}

	switch vendor[0] {
	case postgres:
		portDocker = "5432/tcp"
		env = []string{"POSTGRES_DB=" + ddb.conf.db.Name, "POSTGRES_USER=" + ddb.conf.db.User,
			"POSTGRES_PASSWORD=" + ddb.conf.db.Password}
	case mysql:
		portDocker = "3306/tcp"
		env = []string{"MYSQL_DATABASE=" + ddb.conf.db.Name, "MYSQL_USER=" + ddb.conf.db.User,
			"MYSQL_ROOT_PASSWORD=" + ddb.conf.db.Password,
			"MYSQL_PASSWORD=" + ddb.conf.db.Password}
	default:
		portDocker = ddb.conf.standardDBPort
		env = ddb.conf.envDocker
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			portDocker: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: string(ddb.conf.actualPort),
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}

	containerName := ddb.conf.db.Name
	r, err := ddb.cli.ContainerCreate(ctx, &container.Config{
		Image: specifiedVendor,
		Env:   env,
		ExposedPorts: map[nat.Port]struct{}{
			portDocker: struct{}{},
		},
		Cmd: ddb.conf.cmd,
	}, hostConfig, nil, nil, containerName)
	if err != nil {
		split := strings.Split(err.Error(), `"`)

		if len(split) < 4 {
			return err
		}

		r.ID = split[3]
	}

	ddb.id = r.ID

	err = ddb.setup(ctx)
	if err != nil {
		return err
	}

	return nil
}
