package main

import (
	"fmt"
	"errors"
	"strings"
	"io/ioutil"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type MySql struct {
	Dsn string
	DB  *sql.DB
}

func (this *MySql) prepareDatabase(schemefile string) error {
	bdata, err := ioutil.ReadFile(schemefile)
	if err != nil {
		return err
	}
	commands := strings.Split(string(bdata), ";")
	for _, cmd := range commands {
		if cmd != "" {
			_, err = this.Query(cmd)
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