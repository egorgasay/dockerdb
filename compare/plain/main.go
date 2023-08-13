package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/docker/client"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func getDB(vendor, connStr string) (db *sql.DB, err error) {
	after := time.After(20 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	var errPing error
	for {
		select {
		case <-after:
			if errPing != nil {
				return nil, errPing
			}
			return nil, fmt.Errorf("timeout, Last error:%w", err)
		default:
			db, err = sql.Open(vendor, connStr)
			if db == nil {
				return nil, fmt.Errorf("db is nil %w", err)
			}

			errPing = db.Ping()
			if errPing == nil && err == nil {
				return db, nil
			}

			<-ticker.C
		}
	}
}

func main() {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	_, err = dockerClient.ImagePull(
		context.Background(),
		"postgres",
		types.ImagePullOptions{},
	)
	if err != nil {
		log.Fatal(err)
	}

	containerConfig := &container.Config{
		Image: "postgres",
		Env:   []string{"POSTGRES_PASSWORD=mysecretpassword"},
	}
	containerHostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5432/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5432"}},
		},
	}

	resp, err := dockerClient.ContainerCreate(
		context.Background(),
		containerConfig,
		containerHostConfig,
		nil,
		nil,
		"my-postgres-container",
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	db, err := getDB("postgres",
		"host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	var result string
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	err = dockerClient.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})

	if err != nil {
		log.Fatal(err)
	}
}
