package dockerdb_test

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/egorgasay/dockerdb"
)

func Example() {
	config := dockerdb.CustomDB{
		DB: dockerdb.DB{
			Name:     "admin",
			User:     "admin",
			Password: "admin",
		},
		Port:   "45217",
		Vendor: dockerdb.Postgres15,
	}

	// This will allow you to upload the image to your computer.
	ctx := context.TODO()
	err := dockerdb.Pull(ctx, dockerdb.Postgres15)
	if err != nil {
		log.Fatal(err)
	}

	vdb, err := dockerdb.New(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	var answer string
	err = vdb.DB.QueryRow("SELECT 'db is up'").Scan(&answer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)

	if err = vdb.Stop(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("db is down")
}
