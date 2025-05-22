package db_api
import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"

)
type Session struct {
	ID                string
	UserID            string
	CreatedAt         int64
	CurrentPage       int64
  CurrentPageSearch int64
  Searching         bool
}

func CreateSession(db *sql.DB) (string) {
	createdAt := time.Now().Unix()

  calc := createdAt + 223 + createdAt % 16

  sessionID := "se" + strconv.FormatInt(calc, 10) 

  query := `INSERT INTO sessions (id, UserID, created_at, current_page, current_page_search, searching) VALUES ($1, $2, $3, $4, $5, $6)`
  // FOR searching 
  // 0 = no
  // 1 = yes
	_, err := db.Exec(query, sessionID, nil, createdAt, 1, 1, 0)
  if err != nil{
    println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!111")
    fmt.Println("Error executing query:", err)
  }
	return sessionID
}

func IsLoggedIn(db *sql.DB, sessionID string) bool{
	query := `SELECT COUNT(UserID) FROM sessions WHERE id = $1`
	row := db.QueryRow(query, sessionID)

  var num int 
  err := row.Scan(&num)
	if err != nil {
    log.Fatal(err)
		return false 
  }
  println(num)
  if num == 1 {
    println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
    return true 
  }else{
    return false
  }
}

func GetSession(db *sql.DB, sessionID string) (Session, error) {
	var session Session

	query := `SELECT created_at, current_page, current_page_search FROM sessions WHERE id = $1`
	row := db.QueryRow(query, sessionID)
  println(sessionID)

  err := row.Scan(&session.CreatedAt, &session.CurrentPage, &session.CurrentPageSearch)
	if err != nil {
		return session, err
	}

	return session, nil
}

func UpdateUserSes(db *sql.DB, sessionID string, userID string){
  query := `UPDATE sessions SET userID = $1 WHERE id = $2`

  // Using Exec() for UPDATE query since it doesn't return rows.
  res, err := db.Exec(query, userID, sessionID)
  if err != nil {
    fmt.Println("Error executing query:", err)
  }

  rowsAffected, err := res.RowsAffected()
  if rowsAffected == 0 {
    fmt.Println("No rows updated. Session ID might not exist.")
  } else {
    fmt.Println("Updated", rowsAffected, "row(s).")
  }

}
