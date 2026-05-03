package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8)      NOT NULL DEFAULT "",
    title   VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT         NOT NULL DEFAULT "",
    repeat  VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);
`

func Init() error {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
