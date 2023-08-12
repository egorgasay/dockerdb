package dockerdb

import (
	"github.com/docker/go-connections/nat"
)

func Config() *CustomDB {
	return &CustomDB{
		db: db{
			User:     "test",
			Password: "test",
		},
	}
}

type db struct {
	Name     string
	User     string
	Password string
}

type CustomDB struct {
	db           db
	standardPort string
	vendor       DockerHubName
	vendorName   string

	// Optional if you are using a supported vendor
	actualPort nat.Port
	envDocker  []string
	noSQL      bool
	sqlConnStr string
	pullImage  bool
}

func (c *CustomDB) DBName(name string) *CustomDB {
	c.db.Name = name
	return c
}

func (c *CustomDB) DBUser(user string) *CustomDB {
	c.db.User = user
	return c
}

func (c *CustomDB) DBPassword(password string) *CustomDB {
	c.db.Password = password
	return c
}

func (c *CustomDB) Vendor(vendor DockerHubName) *CustomDB {
	c.vendor = vendor
	return c
}

func (c *CustomDB) StandardPort(port nat.Port) *CustomDB {
	c.standardPort = string(port)
	return c
}

func (c *CustomDB) ActualDBPort(port nat.Port) *CustomDB {
	c.actualPort = port
	return c
}

func (c *CustomDB) DockerEnv(env []string) *CustomDB {
	c.envDocker = env
	return c
}

func (c *CustomDB) SQL() *CustomDB {
	c.noSQL = false
	return c
}

func (c *CustomDB) NoSQL() *CustomDB {
	c.noSQL = true
	return c
}

func (c *CustomDB) SQLConnStr(connString string) *CustomDB {
	c.sqlConnStr = connString
	return c
}

func (c *CustomDB) PullImage() *CustomDB {
	c.pullImage = true
	return c
}

func (c *CustomDB) Build() CustomDB {
	return *c
}
