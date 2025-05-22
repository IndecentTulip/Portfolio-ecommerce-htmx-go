package db_api

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

var postgre_connection = struct {
		user     string
		//password string
		host     string
		port     string
		dbname   string
	}{
		user:     "postgres",  // Replace with actual username
		//password: "password",  // Replace with actual password
		host:     "localhost", // Replace with actual host if necessary
		port:     "5432",      // Default PostgreSQL port
		dbname:   "ecommerce",
	}
//	" password=" + postgre_connection.password +
var connection_str = "user=" + postgre_connection.user +
	" dbname=" + postgre_connection.dbname +
	" host=" + postgre_connection.host +
	" port=" + postgre_connection.port +
	" sslmode=disable"

func CreateDB() *sql.DB{

  db, err := sql.Open("sqlite3", "../../database/products.db")
  if err != nil{
    log.Fatal(err)
  }

	log.Println("Database created successfully.")
  return db
}
func ConnectToDB() *sql.DB{

	db, err := sql.Open("postgres", connection_str)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}
	fmt.Println("Successfully connected to the database!")

	return db
}


