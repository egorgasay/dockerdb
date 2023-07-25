package dockerdb

import "github.com/docker/go-connections/nat"

func Config() *CustomDB {
	return &CustomDB{}
}

type DB struct {
	Name     string
	User     string
	Password string
}

type CustomDB struct {
	DB         DB
	Port       string
	Vendor     DockerHubName
	vendorName string

	// Optional if you are using a supported vendor
	PortDocker nat.Port
	EnvDocker  []string
	NoSQL      bool
	SQLConnStr string
}

func (c *CustomDB) SetName(name string) *CustomDB {
	c.DB.Name = name
	return c
}

func (c *CustomDB) SetUser(user string) *CustomDB {
	c.DB.User = user
	return c
}

func (c *CustomDB) SetPassword(password string) *CustomDB {
	c.DB.Password = password
	return c
}

func (c *CustomDB) SetVendor(vendor DockerHubName) *CustomDB {
	c.Vendor = vendor
	return c
}

func (c *CustomDB) SetDBPort(port string) *CustomDB {
	c.Port = port
	return c
}

func (c *CustomDB) SetExposePort(port nat.Port) *CustomDB {
	c.PortDocker = port
	return c
}

func (c *CustomDB) SetEnv(env []string) *CustomDB {
	c.EnvDocker = env
	return c
}

func (c *CustomDB) SetSQL() *CustomDB {
	c.NoSQL = false
	return c
}

func (c *CustomDB) SetNoSQL() *CustomDB {
	c.NoSQL = true
	return c
}

func (c *CustomDB) SetSQLConnStr(connString string) *CustomDB {
	c.SQLConnStr = connString
	return c
}
