# dockerdb
[![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/mod)](https://pkg.go.dev/golang.org/x/mod)

This repository contains a package for fast database deployment in Docker container.

# Tested Vendors
<ol>
<li>PosgreSQL</li>
<li>MySQL</li>
</ol>

# Why dockerdb?  
  
![изображение](https://user-images.githubusercontent.com/102957432/218540178-a2d56235-076d-400a-a5ac-b83afd49758b.png)

# Usage
Download and install it:
```bash
go get github.com/egorgasay/dockerdb
```

Import it in your code:
```go
import "github.com/egorgasay/dockerdb"
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

# Example 
```go
package main

import (
	"context"
	"fmt"
	"log"
	
	"https://github.com/egorgasay/dockerdb"

    _ "github.com/lib/pq"
)

func main() {
  // Specify the data that is needed to run the database
  config := dockerdb.CustomDB{
    DB: dockerdb.DB{
          Name:     "admin",
          User:     "admin",
          Password: "admin",
      },
      Port:   "45217",
      Vendor: dockerdb.Postgres15,
  }
  
  ctx := context.TODO()
  vdb, err := dockerdb.New(ctx, config)
  if err != nil {
    log.Fatal(err)
  }
  
  // testing db is working
  var answer string
  err = vdb.DB.QueryRow("SELECT 'db is up'").Scan(&answer)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(answer)

  if err = vdb.Stop(ctx); err != nil {
    log.Fatal(err)
  }
}
```

# Example 2 (Unimplemented db)
```go
package main

import (
    "context"
	"fmt"
	"log"
	
	"https://github.com/egorgasay/dockerdb"

    _ "github.com/lib/pq"
)

func main() {
  // Specify the data that is needed to run the database
  config := dockerdb.CustomDB{
    DB: dockerdb.DB{
          Name:     "admin",
          User:     "admin",
          Password: "admin",
      },
      Port:   "45217",
      Vendor: "postgres:10",

    PortDocker: "5432/tcp",
    EnvDocker:  []string{"POSTGRES_DB=" + ddb.conf.DB.Name, "POSTGRES_USER=" + ddb.conf.DB.User,
"POSTGRES_PASSWORD=" + ddb.conf.DB.Password},
  }

	// This will allow you to upload the image to your computer.
	ctx := context.TODO()
	err := dockerdb.Pull(ctx, "postgres:10")
	if err != nil {
		log.Fatal(err)
	}
  
  ctx := context.TODO()
  vdb, err := dockerdb.New(ctx, config)
  if err != nil {
    log.Fatal(err)
  }
  
  // testing db is working
  var answer string
  err = vdb.DB.QueryRow("SELECT 'db is up'").Scan(&answer)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(answer)

  if err = vdb.Stop(ctx); err != nil {
    log.Fatal(err)
  }
}
```
