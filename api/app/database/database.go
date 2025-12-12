package database

import (
	"api/config"

	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

func OpenDB() (db *sqlx.DB, err error) {
	host := config.GetEnvDbHost()
	port := config.GetEnvDbPort()
	name := config.GetEnvDbName()
	user := config.GetEnvDbUser()
	pass := config.GetEnvDbPass()

	// Open DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	if db, err = sqlx.Connect("mysql", dsn); err != nil {
		return db, err
	}

	return db, nil
}

func OpenDBWithTX() (db *sqlx.DB, tx *sqlx.Tx, err error) {
	if db, err = OpenDB(); err != nil {
		return db, tx, err
	}
	if tx, err = db.Beginx(); err != nil {
		return db, tx, err
	}

	return db, tx, nil
}