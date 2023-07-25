module github.com/egorgasay/dockerdb/v2

go 1.19

require (
	github.com/docker/docker v23.0.0+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/go-sql-driver/mysql v1.7.0
	github.com/gocql/gocql v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.7
	github.com/redis/go-redis/v9 v9.0.5
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.11.0

require (
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gotest.tools/v3 v3.4.0 // indirect
)
