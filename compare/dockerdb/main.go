package main

import (
	"context"
	"fmt"
	"github.com/egorgasay/dockerdb/v2"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
)

func main() {
	ctx := context.TODO()
	config := dockerdb.EmptyConfig().DBName("test").DBUser("test").
		DBPassword("test").StandardDBPort("5432").
		Vendor(dockerdb.Postgres15).SQL().PullImage()

	vdb, err := dockerdb.New(ctx, config.Build())
	if err != nil {
		log.Fatal(err)
	}
	defer vdb.Clear(ctx)

	var result string
	err = vdb.SQL().QueryRow("SELECT 'db is up'").Scan(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
