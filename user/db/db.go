package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
}

type UserLogin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func connection() *sql.DB {
	conn, err := sql.Open("sqlite3", "mydatabase.db") // dsn
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

	const createTable string = `CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY NOT NULL, name VARCHAR (255), email VARCHAR (255), password VARCHAR (255))`

	_, err := conn.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertUser(value User) {
	conn := connection()
	timeout, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	conn.PingContext(timeout)

	const insertUser string = `INSERT INTO user (name, email, password) VALUES ($1, $2, $3)`

	_, err := conn.Exec(insertUser, value.Name, value.Mail, value.Password)
	if err != nil {
		log.Fatal(err)
	}
}

func SelectUser(value UserLogin) bool {
	conn := connection()
	const selectUser string = `SELECT id FROM user WHERE name = $1 AND password = $2`
	query, err := conn.Query(selectUser, value.Name, value.Password)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if !query.Next() {
		return false
	}
	fmt.Println(query)
	return true
}
