package dockerdb

import (
	"database/sql"
	"github.com/docker/go-connections/nat"
)

// SQL returns an *sql.DB instance.
func (ddb *VDB) SQL() (db *sql.DB) {
	return ddb.db
}

func (ddb *VDB) GetPort() (port nat.Port) {
	return ddb.conf.actualPort
}
