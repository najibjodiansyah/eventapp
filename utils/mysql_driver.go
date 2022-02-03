package util

import (
	"fmt"

	"database/sql"
	_config "eventapp/config"

	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/gommon/log"
)

var Db *sql.DB

func MysqlDriver(config *_config.AppConfig) *sql.DB {

	uri := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		config.Database.Username,
		config.Database.Password,
		config.Database.Address,
		config.Database.Port,
		config.Database.Name)

	db, err := sql.Open("mysql",uri)

	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}

	return db
}