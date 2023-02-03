package dockerdb

import (
	"context"
	"github.com/docker/docker/api/types"
)

// Run запускает docker контейнер
func (ddb *VDB) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrAlreadyBindPort
		}
	}()

	if err = ddb.cli.ContainerStart(ctx, ddb.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
