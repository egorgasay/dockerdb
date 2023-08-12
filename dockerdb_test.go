package dockerdb_test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/egorgasay/dockerdb/v2"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
	_ "github.com/lib/pq"
	keydb "github.com/redis/go-redis/v9"
)

func TestPostgres15(t *testing.T) {
	config := *dockerdb.
		Config().
		DBName("test").
		DBUser("test").
		DBPassword("test").
		SetDBPort("45217").
		SetVendor(dockerdb.Postgres15).
		SetSQL()

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
	defer vdb.Clear(ctx)

	var answer string
	err = vdb.SQL().QueryRow("SELECT 'db is up'").Scan(&answer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)

	if err = vdb.Stop(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("db is down")
}

func TestKeyDB(t *testing.T) {
	config := dockerdb.Config().
		DBName("mykeydb").DBUser("adm").DBPassword("").
		StandardPort("6379").ActualDBPort("6379/tcp").
		Vendor(dockerdb.KeyDBImage).NoSQL().PullImage().Build()

	ctx := context.TODO()

	vdb, err := dockerdb.New(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer vdb.Clear(ctx)

	var cl *keydb.Client
	var pong string
	for i := 0; i < 10; i++ {
		cl = keydb.NewClient(&keydb.Options{
			Addr:     fmt.Sprintf("%s:%s", "127.0.0.1", "6379"),
			Username: "adm",
			Password: "",
			DB:       0,
		})

		pong, err = cl.Ping(ctx).Result()
		if err == nil {
			break
		}

		log.Println(err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("All is fine: %v\n", pong)

	if err = vdb.Stop(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("db is down")
}

func TestScyllaDB(t *testing.T) {
	cfg := dockerdb.CustomDB{
		db: dockerdb.db{
			Name:     "scylladb",
			User:     "cassandra",
			Password: "cassandra",
		},
		standardPort: "9042",
		vendor:       "scylladb/scylla",
		envDocker: []string{
			"--smp", "1",
			"--authenticator", "PasswordAuthenticator",
			"--broadcast-address", "127.0.0.1",
			"--listen-address", "0.0.0.0",
			"--broadcast-rpc-address", "127.0.0.1",
		},
		actualPort: nat.Port("9042/tcp"),
		noSQL:      true,
	}

	// This will allow you to upload the image to your computer.
	ctx := context.TODO()
	err := dockerdb.Pull(ctx, cfg.vendor)
	if err != nil {
		log.Fatal(err)
	}

	vdb, err := dockerdb.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer vdb.Clear(ctx)

	config := gocql.NewCluster("test")
	config.ProtoVersion = 4
	config.Consistency = gocql.LocalOne
	config.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.db.User, Password: cfg.db.Password,
	}
	config.Hosts = []string{"127.0.0.1"}
	config.Timeout = 5 * time.Second
	config.ConnectTimeout = 5 * time.Second

	for i := 0; i < 30; i++ {
		ses, errSession := gocql.NewSession(*config)
		if errSession == nil {
			ses.Close()
			break
		}
		log.Println(err)

		time.Sleep(2 * time.Second)
	}

	ses, err := gocql.NewSession(*config)
	if err != nil {
		t.Error("session: ", err)
		return
	}
	defer ses.Close()

	if err = ses.Query(`CREATE KEYSPACE IF NOT EXISTS TestKeySpace WITH replication = { 'class': 'NetworkTopologyStrategy', 'replication_factor': '1' } AND durable_writes = TRUE;`).Exec(); err != nil {
		t.Error("CREATE KEYSPACE: ", err)
		return
	}

	config.Keyspace = "TestKeySpace"
	if err != nil {
		t.Error(err)
		return
	}

	var ans string
	err = ses.Query("SELECT key FROM system.local").Scan(&ans)
	if err != nil {
		t.Error("SELECT: ", err)
		return
	}

	if err = vdb.Stop(ctx); err != nil {
		t.Error(err)
		return
	}

	fmt.Println("db is down")
}
