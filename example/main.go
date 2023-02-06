package main

import (
	"context"
	"dockerdb"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config := dockerdb.CustomDB{
		DB: dockerdb.DB{
			Name:     "test",
			User:     "admin",
			Password: "test",
		},
		Port: "35231",
		Vendor: dockerdb.Vendor{
			Name:  dockerdb.Postgres,
			Image: dockerdb.PostgresImage,
		},
	}

	ctx := context.TODO()
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
