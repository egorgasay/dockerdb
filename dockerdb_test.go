package dockerdb_test

import (
	"context"
	"fmt"
	"github.com/egorgasay/dockerdb/v3"
	"github.com/gocql/gocql"

	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	keydb "github.com/redis/go-redis/v9"
	redis "github.com/redis/go-redis/v9"
)

func TestPostgres15(t *testing.T) {
	ctx := context.TODO()
	config := dockerdb.EmptyConfig().DBName("test").DBUser("test").
		DBPassword("test").StandardDBPort("5432").
		Vendor(dockerdb.Postgres15).SQL().PullImage()

	vdb, err := dockerdb.New(ctx, config.Build())
	if err != nil {
		t.Fatal(err)
	}
	defer vdb.Clear(ctx)

	var answer string
	err = vdb.SQL().QueryRow("SELECT 'db is up'").Scan(&answer)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(answer)

	fmt.Println("db is down")
}

func TestPostgresSimple(t *testing.T) {
	ctx := context.TODO()
	vdb, err := dockerdb.New(ctx, dockerdb.PostgresConfig("simple-postgres").Build())
	if err != nil {
		t.Fatal(err)
	}
	defer vdb.Clear(ctx)

	var answer string
	err = vdb.SQL().QueryRow("SELECT 'db is up'").Scan(&answer)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(answer)
}

func TestKeyDB(t *testing.T) {
	var cl *keydb.Client
	var err error
	ctx := context.TODO()

	config := dockerdb.EmptyConfig().
		DBName("mykeydb").StandardDBPort("6379").
		Vendor(dockerdb.KeyDBImage).
		NoSQL(func(conf dockerdb.Config) (stop bool) {
			cl = keydb.NewClient(&keydb.Options{
				Addr: fmt.Sprintf("%s:%s", "127.0.0.1", conf.GetActualPort()),
				DB:   0,
			})

			_, err = cl.Ping(ctx).Result()
			log.Println(err)
			return err == nil
		}, 10, time.Second*2).PullImage()

	vdb, err := dockerdb.New(ctx, config.Build())
	if err != nil {
		t.Fatal(err)
	}
	defer vdb.Clear(ctx)

	fmt.Println("db is down")
}

func TestKeyDBSimple(t *testing.T) {
	var cl *keydb.Client
	var err error
	ctx := context.TODO()

	config := dockerdb.KeyDBConfig("simple-keyDB", func(conf dockerdb.Config) (stop bool) {
		cl = keydb.NewClient(&keydb.Options{
			Addr: fmt.Sprintf("%s:%s", "127.0.0.1", conf.GetActualPort()),
			DB:   0,
		})

		_, err = cl.Ping(ctx).Result()
		log.Println(err)
		return err == nil
	})

	vdb, err := dockerdb.New(ctx, config.Build())
	if err != nil {
		t.Fatal(err)
	}
	defer vdb.Clear(ctx)

	fmt.Println("db is up and everything is fine")
}

func TestRedis(t *testing.T) {
	var cl *keydb.Client
	var err error
	ctx := context.TODO()

	config := dockerdb.EmptyConfig().
		DBName("myredisdb").StandardDBPort("6379").
		Vendor(dockerdb.RedisImage).
		NoSQL(func(conf dockerdb.Config) (stop bool) {
			cl = redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%s", "127.0.0.1", conf.GetActualPort()),
				DB:   0,
			})

			_, err = cl.Ping(ctx).Result()
			log.Println(err)
			return err == nil
		}, 10, time.Second*2).PullImage()

	vdb, err := dockerdb.New(ctx, config.Build())
	if err != nil {
		t.Fatal(err)
	}
	defer vdb.Clear(ctx)

	fmt.Println("db is down")
}

func TestScyllaDB(t *testing.T) {
	ctx := context.TODO()

	// ScyllaDB config prepared
	auth := gocql.PasswordAuthenticator{
		Username: "cassandra", Password: "cassandra",
	}
	config := gocql.NewCluster("test")
	config.ProtoVersion = 4
	config.Consistency = gocql.LocalOne
	config.Authenticator = auth
	config.Timeout = 5 * time.Second
	config.ConnectTimeout = 5 * time.Second

	cfg := dockerdb.EmptyConfig().Vendor(dockerdb.ScyllaDBImage).
		DBName("testScylla").DBUser(auth.Username).DBPassword(auth.Password).
		DockerEnv([]string{
			"--smp", "1",
			"--authenticator", "PasswordAuthenticator",
			"--broadcast-address", "127.0.0.1",
			"--listen-address", "0.0.0.0",
			"--broadcast-rpc-address", "127.0.0.1",
		}).StandardDBPort("9042").NoSQL(func(conf dockerdb.Config) (stop bool) {
		config.Hosts = []string{string("127.0.0.1:" + conf.GetActualPort())}
		_, err := gocql.NewSession(*config)
		if err != nil {
			log.Println(err)
		}
		return err == nil
	}, 30, time.Second*2).PullImage()

	vdb, err := dockerdb.New(ctx, cfg.Build())
	if err != nil {
		t.Fatal(err)
	}
	defer vdb.Clear(ctx)

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

	fmt.Println("db is down")
}
