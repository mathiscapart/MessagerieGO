package db

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type Message struct {
	Dest string

	Exp string

	Message string
}

func connection() *sql.DB {
	conn, err := sql.Open("sqlite3", "db/database.db") // dsn
	if err != nil {
		log.Println(err)
	}

	return conn
}

func Database() {
	conn := connection()
	timeout, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	conn.PingContext(timeout)

	const createTable string = `CREATE TABLE IF NOT EXISTS historique_message (id INTEGER PRIMARY KEY NOT NULL, exp VARCHAR (255), dest VARCHAR(255), mess VARCHAR(255))`

	_, err := conn.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertUser(value Message) {
	conn := connection()
	timeout, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	conn.PingContext(timeout)

	const insertUser string = `INSERT INTO historique_message (exp, dest, mess) VALUES ($1, $2, $3)`

	_, err := conn.Exec(insertUser, value.Exp, value.Dest, value.Message)
	if err != nil {
		log.Fatal(err)
	}
}
