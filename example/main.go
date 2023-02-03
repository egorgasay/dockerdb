package main

import (
	"context"
	"dockerdb"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	i := strconv.Itoa(rand.Int())

	config := dockerdb.CustomDB{
		DB: dockerdb.DB{
			Name:     "test" + i,
			User:     "admin",
			Password: "test",
		},
		Port:   "37053",
		Vendor: "postgres",
	}

	ctx := context.TODO()
	vdb, err := dockerdb.New(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	ctx = context.TODO()
	err = vdb.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var answer string
	err = vdb.DB.QueryRow("SELECT 'db is up'").Scan(&answer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)
}
