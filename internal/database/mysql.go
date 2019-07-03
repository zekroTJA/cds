package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zekroTJA/cds/internal/config"
)

type stmts struct {
	getAccessStat,
	addAccessStat,
	updateAccessStat,
	addRequestLog *sql.Stmt
}

type MySQL struct {
	db *sql.DB

	stmts *stmts
}

func NewMySQL(cfg *config.MySQL) (m *MySQL, err error) {
	m = &MySQL{
		stmts: new(stmts),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.Username, cfg.Password, cfg.Address, cfg.Database)
	if m.db, err = sql.Open("mysql", dsn); err != nil {
		return
	}

	if _, err = m.db.Exec(
		"CREATE TABLE IF NOT EXISTS `accessStats` (" +
			"`uid` int(11) NOT NULL AUTO_INCREMENT," +
			"`fullPath` text NOT NULL," +
			"`fileName` text NOT NULL," +
			"`accesses` bigint(20) NOT NULL DEFAULT '0'," +
			"`lastAccess` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
			"PRIMARY KEY (`uid`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;"); err != nil {
		return
	}

	if _, err = m.db.Exec(
		"CREATE TABLE IF NOT EXISTS `requestLog` (" +
			"`uid` int(11) NOT NULL AUTO_INCREMENT," +
			"`address` text NOT NULL," +
			"`userAgent` text NOT NULL," +
			"`timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
			"`url` text NOT NULL," +
			"`code` int(11) NOT NULL," +
			"PRIMARY KEY (`uid`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=latin1;"); err != nil {
		return
	}

	if m.stmts.addAccessStat, err = m.db.Prepare(
		"INSERT INTO `accessStats`" +
			"(`fullPath`, `fileName`) VALUES (?, ?);"); err != nil {
		return
	}

	if m.stmts.getAccessStat, err = m.db.Prepare(
		"SELECT (`uid`) FROM `accessStats` " +
			"WHERE `fullPath` = ?;"); err != nil {
		return
	}

	if m.stmts.updateAccessStat, err = m.db.Prepare(
		"UPDATE `accessStats` SET `accesses` = `accesses` + 1 " +
			"WHERE `fullPath` = ?;"); err != nil {
		return
	}

	if m.stmts.addRequestLog, err = m.db.Prepare(
		"INSERT INTO `requestLog`" +
			"(`address`, `userAgent`, `url`, `code`)" +
			"VALUES (?, ?, ?, ?)"); err != nil {
		return
	}

	return
}

func (m *MySQL) Close() {
	m.db.Close()
}

func (m *MySQL) RecordAccess(fullPath, fileName, address, userAgent, url string, code int) error {
	err := m.stmts.getAccessStat.QueryRow(fullPath).Scan(new(int))
	if err == sql.ErrNoRows {
		if _, err = m.stmts.addAccessStat.Exec(fullPath, fileName); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if _, err = m.stmts.updateAccessStat.Exec(fullPath); err != nil {
			return err
		}
	}

	_, err = m.stmts.addRequestLog.Exec(address, userAgent, url, code)

	return err
}
