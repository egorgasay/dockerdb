package dockerdb

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Run launches the docker container.
func (ddb *VDB) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrUnknown
		}
	}()

	if err = ddb.cli.ContainerStart(ctx, ddb.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

// Run launches a docker container by ID.
func Run(ctx context.Context, ID string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	ddb := &VDB{
		ID:  ID,
		cli: cli,
	}

	return ddb.Run(ctx)
}

// Pause suspends the docker container.
func (ddb *VDB) Pause(ctx context.Context) (err error) {
	if err = ddb.cli.ContainerPause(ctx, ddb.ID); err != nil {
		return err
	}

	return nil
}

// Unpause resumes the docker container.
func (ddb *VDB) Unpause(ctx context.Context) (err error) {
	if err = ddb.cli.ContainerUnpause(ctx, ddb.ID); err != nil {
		return err
	}

	return nil
}

// Kill kills the docker container.
func (ddb *VDB) Kill(ctx context.Context, signal string) (err error) {
	if err = ddb.cli.ContainerKill(ctx, ddb.ID, signal); err != nil {
		return err
	}

	return nil
}

// Stop stops the docker container.
func (ddb *VDB) Stop(ctx context.Context) (err error) {
	if err = ddb.cli.ContainerStop(ctx, ddb.ID, container.StopOptions{}); err != nil {
		return err
	}

	return nil
}

// Restart stops and starts a container again.
func (ddb *VDB) Restart(ctx context.Context) (err error) {
	if err = ddb.cli.ContainerRestart(ctx, ddb.ID, container.StopOptions{}); err != nil {
		return err
	}

	return nil
}

// Remove removes a container.
func (ddb *VDB) Remove(ctx context.Context) (err error) {
	if err = ddb.cli.ContainerRemove(ctx, ddb.ID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}

// Clear kills and removes a container.
func (ddb *VDB) Clear(ctx context.Context) (err error) {
	if err = ddb.cli.ContainerRemove(ctx, ddb.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	}); err != nil {
		return err
	}

	return nil
}
