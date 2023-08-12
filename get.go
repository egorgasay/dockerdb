package dockerdb

import "github.com/docker/go-connections/nat"

func (c *Config) GetDBUser() string {
	return c.db.User
}

func (c *Config) GetDBPassword() string {
	return c.db.Password
}

func (c *Config) GetDBName() string {
	return c.db.Name
}

func (c *Config) GetEnvDocker() []string {
	return c.envDocker
}

func (c *Config) GetSQLConnStr() string {
	return c.sqlConnStr
}

func (c *Config) GetActualPort() nat.Port {
	return c.actualPort
}

func (c *Config) GetStandardDBPort() nat.Port {
	return c.standardDBPort
}
