package dockerdb

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
)

var supportedVendors = []string{
	"postgres",
	"mysql",
	"mcr.microsoft.com/mssql/server",
}

// Pull pulls an image from net.
// WARNING!! USE IT CAREFULLY! DOWNLOADING SOME DB IMAGES MAY TAKE SOME TIME
func Pull(ctx context.Context, vendor string) error {
	if contains(supportedVendors, vendor) {
		return ErrUnsupportedVendor
	}

	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	pull, err := cli.ImagePull(ctx, vendor, types.ImagePullOptions{})
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
