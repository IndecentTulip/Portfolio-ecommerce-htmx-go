package db_creation

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
  "github.com/google/uuid"

	m "htmxNpython/misc"

	_ "github.com/mattn/go-sqlite3"
)


func CreateProductsTable(db *sql.DB) {
  query := `
    CREATE TABLE IF NOT EXISTS products (
      id TEXT PRIMARY KEY,
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
  checkQuery := `SELECT COUNT(*) FROM products`

  rows,_ := db.Query(checkQuery)
  var num int

  for rows.Next() {
    err := rows.Scan(&num)
    if err != nil {
        log.Fatal(err)
    }
  }
  if num >= 100 {
    println("Products are already inserted")
    return
  }
  
  query := `INSERT INTO products (id, name, price, desc, quantity) VALUES (?, ?, ?, ?, ?)`
  
  for i := 0; i < 100; i++ {
    desc := "lorem Ipsum"
    uniqueID := uuid.New().String()
    if i > 50{ desc = "DEMO DEMO DEMO"}
    _, err := db.Exec(query, uniqueID, "product " + strconv.Itoa(i), 25 + i, desc , 3)
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
  var id string
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
    println(id)
    println(name)
    ProductList = append(ProductList, m.Product{Id: id, Name: name, Price: price, Desc: desc, Quantity: quantity})
  }
  println("product list total:")
  println(total)

  defer rows.Close()

  return ProductList, total

}



func GetProduct(db *sql.DB, prodId string) m.Product {
  query := `SELECT id, name, price, desc, quantity FROM products WHERE id = ?`

  var product m.Product

  rows, err := db.Query(query, prodId)
  if err != nil {
    log.Fatal(err)
  }

  var id string
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
  var id string
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

type Session struct {
	ID                string
	UserID            int
	CreatedAt         int64
	CurrentPage       int64
  CurrentPageSearch int64
  Searching         bool
}

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
  ProductId   string
}

func CreateCartTable(db *sql.DB) {

  query := `
    CREATE TABLE IF NOT EXISTS cart (
      CartId TEXT PRIMARY KEY,
      SessionId TEXT,
      ProductId TEXT,
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
func SelectCart(db *sql.DB, sessionID string) ([]m.CartItem, error) {
    var rowData []m.CartItem

    query := `SELECT 
        c.cartId,
        c.ProductId,
        p.name,
        p.price,
        c.quantity,
        (p.price * c.quantity) AS total
        FROM cart c
        JOIN products p ON c.ProductId = p.id
        WHERE c.SessionId = ?`

    rows, err := db.Query(query, sessionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cartId string
    var productId string
    var name string
    var price int
    var quantity int
    var total int

    for rows.Next() {
        err := rows.Scan(&cartId, &productId, &name, &price, &quantity, &total)
        if err != nil {
            return nil, err
        }

        rowData = append(rowData, 
            m.CartItem{
                Product: m.Product{Id: productId, Name: name, Price: price, Quantity: quantity}, 
                CartID: cartId,
                Total: total,
            })
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return rowData, nil
}
func CountFinalPrice(cartItemsList []m.CartItem) (int) {

  var finalTotal int
  for _,item := range cartItemsList{
    finalTotal += item.Total
  }

  return finalTotal 
}



func SelectCartItem(db *sql.DB, productId string) (m.CartItem, error) {
    var product m.CartItem

    query := `SELECT 
        c.cartId,
        p.name,
        p.price,
        c.quantity,
        (p.price * c.quantity) AS total
        FROM cart c
        JOIN products p ON c.ProductId = p.id
        WHERE c.ProductId = ?`

    row := db.QueryRow(query, productId)

    var cartId string
    var name string
    var price int
    var quantity int
    var total int

    err := row.Scan(&cartId, &name, &price, &quantity, &total)
    if err != nil {
        return product, err
    }

    product.Product = m.Product{Id: productId, Name: name, Price: price, Quantity: quantity}
    product.CartID = cartId

    return product, nil
}

func AddToCart(db *sql.DB, SessionId string, ProductId string){

  query := `INSERT INTO cart
    (CartId, SessionId, ProductId, Quantity)
    VALUES (?, ?, ?, 1)
    ON CONFLICT(SessionId, ProductId) DO UPDATE SET
    Quantity = Quantity + 1;`


  uniqueCartID := uuid.New().String()

  println("!!!!!!!!!!!!!!!!! ADD TO CART !!!!!!!!!!!!!!!!!")
  println(uniqueCartID)
  println(SessionId)
  println(ProductId)

	_, err := db.Exec(query, uniqueCartID, SessionId, ProductId)
  if err != nil {
    log.Fatal(err)
  }


}

func DeleteFromCart(db *sql.DB, cartId string){

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


