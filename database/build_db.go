package build_db

import (
	"database/sql"
	"log"
	crdb "htmxNpython/db_creation"

	_ "github.com/mattn/go-sqlite3"
)

// TODO find a way to execute it without me
// going back and forth chaning name of the func
// maybe make this execute when you pass a param to the main.go or something
func CreateDummyDB(){

  db, err := sql.Open("sqlite3", "../../database/products.db")
  if err != nil{
    log.Fatal(err)
  }

  crdb.DeleteDB()
  crdb.CreateDB()

  // make it drop all the db

  crdb.CreateUsersTable(db)
  crdb.CreateProductsTable(db)
  crdb.CreateTagsTable(db)
  crdb.CreateTagsForProductTable(db)
  crdb.CreateSessionsTable(db)
  crdb.CreateCartTable(db)
  crdb.InsertIntoProducts(db)
  crdb.InsertDefaultTags(db)
  crdb.InsertDefaultTags_for_Products(db)

}


