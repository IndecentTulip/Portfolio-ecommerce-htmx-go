package db_api
import (
	"database/sql"
	"log"

	wc "HtmxReactGolang/web_context"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"

)


func InsertIntoUser(db *sql.DB, user wc.UserContext){
  checkQuery := `SELECT COUNT(*) FROM users WHERE id == $1`

  row := db.QueryRow(checkQuery, user.UserID)
  var num int
  err := row.Scan(&num)
  if err != nil {
    log.Fatal(err)
  }

  if num > 0 {
    println("User already inserted")
    return
  }

  query := `INSERT INTO users (id, name, profileImage) VALUES ($1, $2, $3)`

  _, err = db.Exec(query, user.UserID, user.UserName, user.ProfileImage)
  if err != nil {
      log.Fatal(err)
  }

  fmt.Println("User Data inserted successfully!")

}


func GetUser(db *sql.DB, sessionID string) wc.UserContext{
  query := `SELECT u.name, u.profileImage
  FROM sessions s
  JOIN users u ON s.UserID = u.id
  WHERE s.id = $1;`
  
  row := db.QueryRow(query, sessionID)

  var name, image string

  err := row.Scan(&name, &image)
  if err != nil {
    log.Fatal(err)
  }

  user :=  wc.UserContext{
    UserName: name,
    ProfileImage: image,
  }

  println("INSIDE GET USER")
  println(name)
  
  return user

}


