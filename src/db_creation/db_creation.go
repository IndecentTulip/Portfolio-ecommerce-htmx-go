package db_creation

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	m "htmxNpython/misc"

	_ "github.com/mattn/go-sqlite3"
)


func CreateProductsTable(db *sql.DB) {
  query := `
    CREATE TABLE IF NOT EXISTS products (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT,
      price INTEGER,
      desc TEXT,
      quantity INTEGER
      );
  `
  _, err := db.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Product Table created successfully!")
}

func insertIntoProducts(db *sql.DB) {
    rows, err := db.Query("SELECT id FROM products")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
      var id int
        if err := rows.Scan(&id); err != nil {
          log.Fatal(err)
        }
        if id == 100 {
          log.Printf("Products already created")
          return
        }
    }

    query := `INSERT INTO products (name, price, desc, quantity) VALUES (?, ?, ?, ?)`
    
    for i := 0; i < 100; i++ {
      desc := "lorem Ipsum"
      if i > 50{ desc = "DEMO DEMO DEMO"}
      _, err := db.Exec(query, "product " + strconv.Itoa(i), 25 + i, desc , 3)
      if err != nil {
          log.Fatal(err)
      }
    }

    fmt.Println("Data inserted successfully!")
}

func GetProductsList(db *sql.DB,  offset int) ([]m.Product, int) {
  query := `
      SELECT
      COUNT(*) OVER() AS total,
      id, name, price, desc, quantity 
      FROM products 
      LIMIT 10
      OFFSET ?`
      
  var ProductList []m.Product

  rows, err := db.Query(query, offset)
  if err != nil {
    log.Fatal(err)
  }

  var total int
  var id int
  var name string
  var price int
  var desc string
  var quantity int

  for rows.Next() {
    err := rows.Scan(&total, &id, &name, &price, &desc, &quantity)
    if err != nil {
        log.Fatal(err)
    }
    
    println("product list: adding:")
    println(name)
    ProductList = append(ProductList, m.Product{Id: id, Name: name, Price: price, Desc: desc, Quantity: quantity})
  }
  println("product list total:")
  println(total)

  defer rows.Close()

  return ProductList, total

}



func GetProduct(db *sql.DB, prodId int) m.Product {
  query := `SELECT id, name, price, desc, quantity FROM products WHERE id = ?`

  var product m.Product

  rows, err := db.Query(query, prodId)
  if err != nil {
    log.Fatal(err)
  }

  var id int
  var name string
  var price int
  var desc string
  var quantity int

  for rows.Next() {
    err := rows.Scan(&id, &name, &price, &desc, &quantity)
    if err != nil {
        log.Fatal(err)
    }
    
    product = m.Product{Id: id, Name: name, Price: price, Desc: desc, Quantity: quantity}
  }

    defer rows.Close()

  return product

}

func GetProductSearch(db *sql.DB, term string, offset int) ([]m.Product, int) {
  terms := strings.Split(term, " ")
  query := `
      SELECT
      COUNT(*) OVER() AS total,
      id, name, price, desc, quantity 
      FROM products 
      WHERE ` + helpterBuildWhereClause(len(terms)) + `
      COLLATE NOCASE
      LIMIT 10
      OFFSET ?`
      

  params := make([]interface{}, 0, len(terms)*2)
  for _, term := range terms {
      term = "%" + term + "%"
      params = append(params, term, term)
  }
  params = append(params, offset)

  var searchProductList []m.Product

  rows, err := db.Query(query, params...)
  if err != nil {
    log.Fatal(err)
  }

  var total int
  var id int
  var name string
  var price int
  var desc string
  var quantity int

  for rows.Next() {
    err := rows.Scan(&total, &id, &name, &price, &desc, &quantity)
    if err != nil {
        log.Fatal(err)
    }
    
    println("search: adding:")
    println(name)
    searchProductList = append(searchProductList, m.Product{Id: id, Name: name, Price: price, Desc: desc, Quantity: quantity})
  }
  println("search total:")
  println(total)

  defer rows.Close()

  return searchProductList, total

}

func helpterBuildWhereClause(termCount int) string {
    var clauses []string
    for i := 0; i < termCount; i++ {
        clauses = append(clauses, "(name LIKE ? OR desc LIKE ?)")
    }
    return strings.Join(clauses, " AND ")
}


//func queryData(db *sql.DB) {
//    rows, err := db.Query("SELECT id, name, age FROM users")
//    if err != nil {
//        log.Fatal(err)
//    }
//    defer rows.Close()
//
//    for rows.Next() {
//        var id int
//        var name string
//        var age int
//        if err := rows.Scan(&id, &name, &age); err != nil {
//            log.Fatal(err)
//        }
//        fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
//    }
//}
//
//func updateData(db *sql.DB) {
//    query := `UPDATE users SET age = ? WHERE name = ?`
//    _, err := db.Exec(query, 30, "Alice")
//    if err != nil {
//        log.Fatal(err)
//    }
//    fmt.Println("Data updated successfully!")
//}
//
//func deleteData(db *sql.DB) {
//    query := `DELETE FROM users WHERE name = ?`
//    _, err := db.Exec(query, "Alice")
//    if err != nil {
//        log.Fatal(err)
//    fmt.Println("Data deleted successfully!")
//}

type Session struct {
	ID                string
	UserID            int
	CreatedAt         int64
	CurrentPage       int64
  CurrentPageSearch int64
  Searching         bool
}

// user_id INTEGER,
func CreateSessionsTable(db *sql.DB) {
  query := `
    CREATE TABLE IF NOT EXISTS sessions (
      id TEXT PRIMARY KEY,
      created_at INTEGER,
      current_page INTEGER,
      current_page_search INTEGER, 
      searching INTEGER
    );
  `
  _, err := db.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Session Table created successfully!")
}

func CreateSession(db *sql.DB) (string, error) {
	createdAt := time.Now().Unix()

  calc := createdAt + 223 + createdAt % 16

  sessionID := "se" + strconv.FormatInt(calc, 10) 

  query := `INSERT INTO sessions (id, created_at, current_page, current_page_search, searching) VALUES (?, ?, ?, ?, ?)`
  // FOR searching 
  // 0 = no
  // 1 = yes
	_, errexec := db.Exec(query, sessionID, createdAt, 1, 1, 0)
	return sessionID, errexec
}

func UpdatePageNumSes(db *sql.DB, sessionID string, num int) error {
  query := `UPDATE sessions SET current_page = ? WHERE id = ?`

  // Using Exec() for UPDATE query since it doesn't return rows.
  res, err := db.Exec(query, num, sessionID)
  if err != nil {
    fmt.Println("Error executing query:", err)
    return err
  }

  rowsAffected, err := res.RowsAffected()
  if rowsAffected == 0 {
    fmt.Println("No rows updated. Session ID might not exist.")
  } else {
    fmt.Println("Updated", rowsAffected, "row(s).")
  }

  return nil
}

func UpdatePageSearchNumSes(db *sql.DB, sessionID string, num int) error {
  query := `UPDATE sessions SET current_page_search = ? WHERE id = ?`

  // Using Exec() for UPDATE query since it doesn't return rows.
  res, err := db.Exec(query, num, sessionID)
  if err != nil {
    fmt.Println("Error executing query:", err)
    return err
  }

  rowsAffected, err := res.RowsAffected()
  if rowsAffected == 0 {
    fmt.Println("No rows updated. Session ID might not exist.")
  } else {
    fmt.Println("Updated", rowsAffected, "row(s).")
  }

  return nil
}
func UpdateSearchingStatus(db *sql.DB, sessionID string, status bool) error {
  query := `UPDATE sessions SET searching = ? WHERE id = ?`

  var statusInt int
  if status {
    statusInt = 1
  }else{
    statusInt = 0
  }
  // Using Exec() for UPDATE query since it doesn't return rows.
  res, err := db.Exec(query, statusInt, sessionID)
  if err != nil {
    fmt.Println("Error executing query:", err)
    return err
  }

  rowsAffected, err := res.RowsAffected()
  if rowsAffected == 0 {
    fmt.Println("No rows updated. Session ID might not exist.")
  } else {
    fmt.Println("Updated", rowsAffected, "row(s).")
  }

  return nil
}



func GetSession(db *sql.DB, sessionID string) (Session, error) {
	var session Session

	query := `SELECT created_at, current_page, current_page_search FROM sessions WHERE id = ?`
	row := db.QueryRow(query, sessionID)
  println(sessionID)

  err := row.Scan(&session.CreatedAt, &session.CurrentPage, &session.CurrentPageSearch)
	if err != nil {
		return session, err
	}

	return session, nil
}

type Cart struct{
  SessionId   string
  ProductId   int
}

func CreateCartTable(db *sql.DB) {

  query := `
    CREATE TABLE IF NOT EXISTS cart (
      CartId INTEGER PRIMARY KEY AUTOINCREMENT,
      SessionId TEXT,
      ProductId INTEGER,
      Quantity INTEGER DEFAULT 1,
      FOREIGN KEY (SessionId) REFERENCES sessions(id),
      FOREIGN KEY (ProductId) REFERENCES products(id),
      UNIQUE (SessionId, ProductId) 
    );
  `
  _, err := db.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Cart Table created successfully!")
}

func SelectCart(db *sql.DB, sessionID string) ([]m.Product, error) {

  var rowData []m.Product

  query := `SELECT 
    c.cartId,
    p.name,
    p.price,
    c.quantity,
    (p.price * c.quantity) AS total
    FROM cart c
    JOIN products p ON c.ProductId = p.id
    WHERE c.SessionId = ?`

    rows, err := db.Query(query, sessionID)
    if err != nil {
      log.Fatal(err)
    }
    //var product Product
    var cartId int
    var name string
    var price int
    var quantity int
    var total int

    for rows.Next() {
      err := rows.Scan(&cartId, &name, &price, &quantity, &total)
      if err != nil {
          log.Fatal(err)
      }
      
      println("adding:")
      println(name)
      rowData = append(rowData, m.Product{Id: cartId, Name: name, Price: price, Quantity: quantity})
      println(rowData[0].Price)
    }

    defer rows.Close()
  return rowData,nil
}

func SelectCartItem(db *sql.DB, productId int) (m.Product, error) {

  var product m.Product

  query := `SELECT 
    c.cartId,
    p.name,
    p.price,
    c.quantity,
    (p.price * c.quantity) AS total
    FROM cart c
    JOIN products p ON c.ProductId = p.id
    WHERE c.ProductId = ?
    `

    row, err := db.Query(query, productId)
    if err != nil {
      log.Fatal(err)
    }
    //var product Product
    var cartId int
    var name string
    var price int
    var quantity int
    var total int

    for row.Next() {
      err := row.Scan(&cartId, &name, &price, &quantity, &total)
      if err != nil {
          log.Fatal(err)
      }
      
      product = m.Product{Id: cartId, Name: name, Price: price, Quantity: quantity}
    }

    defer row.Close()
  return product,nil
}


func AddToCart(db *sql.DB, SessionId string, ProductId int){

  query := `INSERT INTO cart
    (SessionId, ProductId, Quantity)
    VALUES (?, ?, 1)
    ON CONFLICT(SessionId, ProductId) DO UPDATE SET
    Quantity = Quantity + 1;`

	_, err := db.Exec(query, SessionId, ProductId)
  if err != nil {
    log.Fatal(err)
  }


}

func DeleteFromCart(db *sql.DB, cartId int){

  query := `DELETE FROM cart WHERE CartId = ?`

	_, err := db.Exec(query, cartId)
  if err != nil {
    log.Fatal(err)
  }
}


func CreateDB() *sql.DB{

  db, err := sql.Open("sqlite3", "../database/products.db")
  if err != nil{
    log.Fatal(err)
  }

  CreateProductsTable(db)
  CreateSessionsTable(db)
  CreateCartTable(db)
  insertIntoProducts(db)

  

  return db
}


