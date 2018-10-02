package main

import (
	"fmt"
	"errors"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var dbScheme = `
CREATE TABLE IF NOT EXISTS accessStats (
	fullPath text NOT NULL,
	fileName text NOT NULL,
	accesses bigint(20) NOT NULL DEFAULT '0',
	lastAccess timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
  
CREATE TABLE IF NOT EXISTS requestLog (
	address text NOT NULL,
	userAgent text NOT NULL,
	timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	url text NOT NULL,
	code int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
`

type MySql struct {
	Dsn string
	DB  *sql.DB
}

func (this *MySql) prepareDatabase(schemefile string) error {
	commands := strings.Split(dbScheme, ";")
	for _, cmd := range commands {
		if cmd != "" {
			_, err := this.Query(cmd)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewMySql(address, username, password, database, schemefile string) (*MySql, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, address, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	mysql := &MySql{ dsn, db }
	_, err = mysql.Query("SHOW TABLES;")
	if err != nil {
		return nil, err
	}
	err = mysql.prepareDatabase(schemefile)
	if err != nil {
		return nil, err
	}
	return mysql, nil
}

func (this *MySql) Close() {
	if this == nil {
		return
	}
	this.DB.Close()
}

func (this *MySql) Query(statement string, values ...interface{}) (*sql.Rows, error) {
	if this == nil {
		return nil, errors.New("nullptr")
	}
	stm, err := this.DB.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	return stm.Query(values...)
}