package dockerdb

import (
	"fmt"
)

// Build builds connection string by CustomDB config.
func Build(conf CustomDB) (connStr string, err error) {
	switch conf.vendorName {
	case "postgres":
		return fmt.Sprintf(
			"host=localhost user=%s password='%s' dbname=%s port=%s sslmode=disable",
			conf.db.User, conf.db.Password, conf.db.Name, conf.standardPort), nil
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(127.0.0.1:%s)/%s",
			conf.db.User, conf.db.Password, conf.standardPort, conf.db.Name), nil
	default:
		return "", ErrUnsupported
	}
}
