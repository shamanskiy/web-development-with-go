package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db, err := sql.Open("pgx", `
	host=localhost 
	port=5432  
	user=baloo 
	password=junglebook 
	dbname=lenslocked 
	sslmode=disable`)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		amount INT,
		description TEXT
	);`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created.")

	name := "New Calhoun"
	email := "new@calhoun.io"
	row := db.QueryRow(`
  INSERT INTO users(name, email)
  VALUES($1, $2) RETURNING id;`, name, email)

	var id int
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("New user with id %d inserted\n", id)
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (pc PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s  user=%s password=%s dbname=%s sslmode=%s", pc.Host, pc.Port, pc.User, pc.Password, pc.Database, pc.SSLMode)
}
