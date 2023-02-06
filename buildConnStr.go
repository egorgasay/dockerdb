package dockerdb

import "fmt"

func Build(conf CustomDB) (connStr string) {
	switch conf.Vendor.Name {
	case "postgres":
		return fmt.Sprintf(
			"host=localhost user=%s password='%s' dbname=%s port=%s sslmode=disable",
			conf.DB.User, conf.DB.Password, conf.DB.Name, conf.Port)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(127.0.0.1:%s)/%s",
			conf.DB.User, conf.DB.Password, conf.Port, conf.DB.Name)
	//case "mssql":
	//	return fmt.Sprintf(
	//		"Server=127.0.0.1,%s;Database=msdb;User Id=sa;Password=%s;",
	//		conf.Port, conf.DB.Password)
	default:
		return ""
	}
}
