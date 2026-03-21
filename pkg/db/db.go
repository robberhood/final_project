package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL,
		comment TEXT,
		repeat VARCHAR(128) NOT NULL);
		CREATE INDEX idx_tasks_date ON scheduler(date);`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	if install {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}
