package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/robberhood/final_project/config"
	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT,
		repeat VARCHAR(128) NOT NULL DEFAULT "");
		CREATE INDEX idx_tasks_date ON scheduler(date);`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
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

func AddTask(task *config.Task) (int64, error) {

	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err

}

func Tasks(limit int) ([]*config.Task, error) {
	query := `SELECT id,date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	var tasks []*config.Task

	rows, err := db.Query(query, limit)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {

		var t config.Task

		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Title,
			&t.Comment,
			&t.Repeat,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &t)
	}

	if tasks == nil {
		tasks = []*config.Task{}
	}
	return tasks, nil
}
func GetTask(id string) (*config.Task, error) {
	var t config.Task
	query := `SELECT id,date, title, comment, repeat FROM scheduler WHERE id=?`
	err := db.QueryRow(query, id).Scan(
		&t.ID,
		&t.Date,
		&t.Title,
		&t.Comment,
		&t.Repeat,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &t, nil
}

func UpdateTask(task *config.Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id=?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}

func UpdateDate(id string, date string) error {
	query := `UPDATE scheduler SET date=? WHERE id=?`
	res, err := db.Exec(query, date, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func TasksByDate(date string) ([]*config.Task, error) {
	query := `SELECT * FROM scheduler WHERE date=? LIMIT ?`
	var tasks []*config.Task

	rows, err := db.Query(query, date, 50)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {

		var t config.Task

		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Title,
			&t.Comment,
			&t.Repeat,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &t)
	}

	if tasks == nil {
		tasks = []*config.Task{}
	}
	return tasks, nil
}

func TasksByText(text string) ([]*config.Task, error) {
	query := `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
	var tasks []*config.Task

	rows, err := db.Query(query, "%"+text+"%", "%"+text+"%", 50)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {

		var t config.Task

		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Title,
			&t.Comment,
			&t.Repeat,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &t)
	}

	if tasks == nil {
		tasks = []*config.Task{}
	}
	return tasks, nil
}
