package dockerdb

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"log"
)

var supportedVendors = []string{
	"postgres",
	"mysql",
	"mcr.microsoft.com/mssql/server",
}

// Pull pulls an image from net.
// WARNING!! USE IT CAREFULLY! DOWNLOADING SOME DB IMAGES MAY TAKE SOME TIME
func (ddb *VDB) Pull(ctx context.Context, vendor string) error {
	if !ddb.contains(supportedVendors, vendor) {
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
