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
			conf.DB.User, conf.DB.Password, conf.DB.Name, conf.Port), nil
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(127.0.0.1:%s)/%s",
			conf.DB.User, conf.DB.Password, conf.Port, conf.DB.Name), nil
	default:
		return "", ErrUnsupported
	}
}
