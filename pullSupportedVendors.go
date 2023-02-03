package dockerdb

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"log"
)

var SupportedVendors = []string{
	"postgres",
	"mysql",
}

func (ddb *VDB) pull(ctx context.Context, vendor string) error {
	if !ddb.contains(SupportedVendors, vendor) {
		return ErrUnsupportedVendor
	}

	pull, err := ddb.cli.ImagePull(ctx, vendor, types.ImagePullOptions{})
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

func (ddb *VDB) contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
