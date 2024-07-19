package driver

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgconn"

	_ "github.com/jackc/pgx/v4"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func CreateDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 dbname=hackhive user=postgres password=raptor3796")
	if err != nil {
		log.Fatal("Could not create database")
		return nil, err
	}
	defer db.Close()

	return db, nil
}

func DisplayRows(db *sql.DB) error {
	rows, err := db.Query("select id, username, password, email, phone from login_details")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	var username, password, email, phone string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &username, &password, &email, &phone)
		if err != nil {
			log.Println("Error scanning", err)
			return err
		}
		fmt.Println("Record is", id, username, password, email, phone)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning errors", err)
	}

	fmt.Println("-----------------------------------------")

	return nil
}
