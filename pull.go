package dockerdb

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
)

// Pull pulls an image from net.
// WARNING!! USE IT CAREFULLY! DOWNLOADING SOME db IMAGES MAY TAKE SOME TIME.
// Tested with PostgreSQL, MySQL, MS SQL.
func Pull(ctx context.Context, image DockerHubName) error {
	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	pull, err := cli.ImagePull(ctx, string(image), types.ImagePullOptions{})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(pull)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	err = pull.Close()
	if err != nil {
		return err
	}

	return nil
}
