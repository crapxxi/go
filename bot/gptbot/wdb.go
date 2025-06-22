package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func addData(data User, db *sql.DB) error {
	_, err := db.Exec("insert into users (telegram_id, name, level, target_lang, progress) values ($1, $2, $3, $4, $5)", data.telegram_id, data.name, data.level, data.target_lang, data.progress)
	return err
}
func updateData(tid int64, data string, target string, db *sql.DB) error {
	if target == "level" {
		_, err := db.Exec("update users set level = $1 where telegram_id = $2", data, tid)
		return err
	} else if target == "progress" {
		lastprogress, err := getProgress(db, tid)
		if err != nil {
			return err
		}
		_, err = db.Exec("update users set progress = $1 where telegram_id = $2", lastprogress+" "+data+" ", tid)
		return err
	}
	_, err := db.Exec("update users set target_lang = $1 where telegram_id = $2", data, tid)
	return err
}
func getProgress(db *sql.DB, tid int64) (string, error) {
	row := db.QueryRow("select COALESCE(progress, '') from users where telegram_id = $1", tid)
	var progress string
	err := row.Scan(&progress)
	if err != nil {
		return "", err
	}
	return progress, nil
}
func getLevel(db *sql.DB, tid int64) (string, error) {
	row := db.QueryRow("select level from users where telegram_id = $1", tid)
	var level string
	err := row.Scan(&level)
	if err != nil {
		return "", err
	}
	return level, nil
}

func getLanguage(db *sql.DB, tid int64) (string, error) {
	row := db.QueryRow("select target_lang from users where telegram_id = $1", tid)
	var lang string
	err := row.Scan(&lang)
	if err != nil {
		return "", err
	}
	return lang, nil
}

func deleteprogress(tid int64, db *sql.DB) error {
	_, err := db.Exec("update users set progress = null where telegram_id = $1", tid)
	return err
}

func wdb() (*sql.DB, error) {
	connStr := "user=postgres password=admin dbname=langptbot sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	return db, err
}
