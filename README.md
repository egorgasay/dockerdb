# dockerdb
[![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/mod)](https://pkg.go.dev/golang.org/x/mod)

This repository contains a package for fast database deployment in Docker container.

# Supported Vendors
<ol>
<li>PosgreSQL</li>
<li>MySQL</li>
</ol>
More vendors soon...

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
import (
  "context"
  "https://github.com/egorgasay/dockerdb"
  "fmt"
  "log"
  "strconv"
  "time"

  _ "github.com/lib/pq"
)

func main() {
  // Specify the data that is needed to run the database
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

