package dockerdb

import (
	"errors"
	"github.com/docker/go-connections/nat"
	"time"
)

func EmptyConfig() *Config {
	return &Config{
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

type Config struct {
	db             db
	standardDBPort nat.Port
	vendor         DockerHubName
	vendorName     string

	// Optional
	actualPort  nat.Port
	envDocker   []string
	sqlConnStr  string
	noSQL       bool
	checkWakeUp checkWakeUp
	pullImage   bool
}

type checkWakeUp struct {
	fn        func(conf Config) (stop bool)
	sleepTime time.Duration
	tries     int
}

// DBName sets the name of the database.
func (c *Config) DBName(name string) *Config {
	c.db.Name = name
	return c
}

// DBUser sets the user of the database.
func (c *Config) DBUser(user string) *Config {
	c.db.User = user
	return c
}

// DBPassword sets the password of the database.
func (c *Config) DBPassword(password string) *Config {
	c.db.Password = password
	return c
}

// Vendor sets the vendor of the database.
func (c *Config) Vendor(vendor DockerHubName) *Config {
	c.vendor = vendor
	return c
}

// ActualPort allows you to set the actual port for the database.
// Random unused port is used by default.
func (c *Config) ActualPort(port nat.Port) *Config {
	c.actualPort = port
	return c
}

// StandardDBPort represents the standard port of the database which can be used to connect to it.
func (c *Config) StandardDBPort(port nat.Port) *Config {
	c.standardDBPort = port
	return c
}

// DockerEnv sets the environment variables for docker.
func (c *Config) DockerEnv(env []string) *Config {
	c.envDocker = env
	return c
}

// SQL sets db kind to SQL.
func (c *Config) SQL() *Config {
	c.noSQL = false
	return c
}

// NoSQL sets db kind to NoSQL.
func (c *Config) NoSQL(checkWakeUp func(conf Config) (stop bool), tries int, sleepTime time.Duration) *Config {
	c.checkWakeUp.fn = checkWakeUp
	c.checkWakeUp.tries = tries
	c.checkWakeUp.sleepTime = sleepTime
	c.noSQL = true
	return c
}

// SQLConnStr sets the SQL connection string. (Use it only for unsupported databases).
func (c *Config) SQLConnStr(connString string) *Config {
	c.sqlConnStr = connString
	return c
}

// PullImage pulls the vendor image.
func (c *Config) PullImage() *Config {
	c.pullImage = true
	return c
}

// Build builds the config. After building, the config can be used and can't be changed.
func (c *Config) Build() Config {
	return *c
}

func validate(c Config) error {
	if c.vendor == "" {
		return errors.New("vendor must be not empty")
	}

	if c.standardDBPort == "" {
		return errors.New("port must be not empty")
	}

	if c.db.Name == "" {
		return errors.New("db name must be not empty")
	}

	return nil
}
