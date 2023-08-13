# dockerdb
[![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/mod)](https://pkg.go.dev/golang.org/x/mod)

This repository contains a package for fast database deployment in Docker container.

# Why dockerdb?  
  
![изображение](https://user-images.githubusercontent.com/102957432/218540178-a2d56235-076d-400a-a5ac-b83afd49758b.png)

# Usage
Download and install it:
```bash
go get github.com/egorgasay/dockerdb/v3
```

Import it in your code:
```go
import "github.com/egorgasay/dockerdb/v3"
```

The first launch should look like this:
```go
vdb, err := dockerdb.New(ctx, config)
if err != nil {
  log.Fatal(err)
}
```

If the database was turned off, then you can turn it on using:
```go
err := vdb.Run(ctx)
if err != nil {
  log.Fatal(err)
}
```

# SQL DB Example 
```go
package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/egorgasay/dockerdb/v3"

	_ "github.com/lib/pq"
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

	if err = vdb.Stop(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("db is down")
}
```

# NoSQL DB Example
```go
package main

import (
    "context"
	"fmt"
	"log"
	
	"github.com/egorgasay/dockerdb/v3"
	
	redis "github.com/redis/go-redis/v9"
)

func main() {
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
		log.Fatal(err)
	}
	defer vdb.Clear(ctx)

	fmt.Println("db is up")

	if err = vdb.Stop(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("db is down")
}
```
