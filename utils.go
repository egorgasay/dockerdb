package dockerdb

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"net"
	"strconv"
	"time"
)

func SetMaxWaitTime(sec time.Duration) {
	maxWaitTime = sec
}

func getFreePort() (nat.Port, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "0", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "0", err
	}

	defer l.Close()
	port := strconv.Itoa(l.Addr().(*net.TCPAddr).Port)

	return nat.Port(port), nil
}

// buildConnStr builds connection string by CustomDB config.
func buildConnStr(conf Config) (connStr string, err error) {
	switch conf.vendorName {
	case "postgres":
		return fmt.Sprintf(
			"host=localhost user=%s password='%s' dbname=%s port=%s sslmode=disable",
			conf.db.User, conf.db.Password, conf.db.Name, conf.actualPort), nil
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(127.0.0.1:%s)/%s",
			conf.db.User, conf.db.Password, conf.actualPort, conf.db.Name), nil
	default:
		return "", ErrUnsupported
	}
}
