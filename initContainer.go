package dockerdb

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// init initializes the docker container with the selected database
func (ddb *VDB) init(ctx context.Context) error {
	var env []string
	var portDocker nat.Port

	if ddb.conf.Port == "" {
		return errors.New("port must be not empty")
	}

	switch ddb.conf.Vendor.Name {
	case "postgres":
		portDocker = "5432/tcp"
		env = []string{"POSTGRES_DB=" + ddb.conf.DB.Name, "POSTGRES_USER=" + ddb.conf.DB.User,
			"POSTGRES_PASSWORD=" + ddb.conf.DB.Password}
	case "mysql":
		portDocker = "3306/tcp"
		env = []string{"MYSQL_DATABASE=" + ddb.conf.DB.Name, "MYSQL_USER=" + ddb.conf.DB.User,
			"MYSQL_ROOT_PASSWORD=" + ddb.conf.DB.Password,
			"MYSQL_PASSWORD=" + ddb.conf.DB.Password}
		//case "mssql":
		//	portDocker = "1433/tcp"
		//	env = []string{"ACCEPT_EULA=Y", "MSSQL_SA_PASSWORD=" + ddb.conf.DB.Password}
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
		Image: ddb.conf.Vendor.Image,
		Env:   env,
		//Cmd:   []string{"sh", "-c", "while true; do sleep 10; done"},
	}, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	ddb.ID = r.ID

	err = ddb.setup(ctx)
	if err != nil {
		return err
	}

	return nil
}
