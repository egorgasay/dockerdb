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

	if ddb.conf.Port == "" {
		return errors.New("port must be not empty")
	}

	vendor := strings.Split(ddb.conf.Vendor, ":")
	if len(vendor) == 0 {
		return errors.New("vendor must be not empty")
	}

	switch vendor[0] {
	case "postgres":
		portDocker = "5432/tcp"
		env = []string{"POSTGRES_DB=" + ddb.conf.DB.Name, "POSTGRES_USER=" + ddb.conf.DB.User,
			"POSTGRES_PASSWORD=" + ddb.conf.DB.Password}
	case "mysql":
		portDocker = "3306/tcp"
		env = []string{"MYSQL_DATABASE=" + ddb.conf.DB.Name, "MYSQL_USER=" + ddb.conf.DB.User,
			"MYSQL_ROOT_PASSWORD=" + ddb.conf.DB.Password,
			"MYSQL_PASSWORD=" + ddb.conf.DB.Password}
	default:
		portDocker = ddb.conf.PortDocker
		env = ddb.conf.EnvDocker
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			portDocker: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: ddb.conf.Port,
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}

	containerName := ddb.conf.DB.Name
	r, err := ddb.cli.ContainerCreate(ctx, &container.Config{
		Image: ddb.conf.Vendor,
		Env:   env,
	}, hostConfig, nil, nil, containerName)
	if err != nil {
		split := strings.Split(err.Error(), `"`)

		if len(split) < 4 {
			return err
		}

		r.ID = split[3]
	}

	ddb.ID = r.ID

	err = ddb.setup(ctx)
	if err != nil {
		return err
	}

	return nil
}
