package db_creation

import (
	"math/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	wc "htmxNpython/web_context"

	_ "github.com/mattn/go-sqlite3"
)
var productTags = [60]string{
  "electronics", "clothing", "home-appliances", "books", "furniture",
  "sports", "toys", "automotive", "new-arrival", "best-seller",
  "limited-edition", "exclusive", "eco-friendly", "handmade", "premium",
  "affordable", "summer-sale", "winter-collection", "black-friday",
  "holiday-special", "back-to-school", "sale", "clearance", "discounted",
  "buy-one-get-one", "free-shipping", "flash-deal", "new", "refurbished",
  "used", "vintage", "open-box", "kids", "adults", "women", "men", "seniors",
  "unisex", "nike", "samsung", "apple", "sony", "adidas", "puma", "lego",
  "durable", "lightweight", "waterproof", "ergonomic", "energy-efficient",
  "fast-charging", "multi-purpose", "high-performance", "red", "blue", "black",
  "white", "modern", "minimalistic", "classic",
}
var productTagsMap = map[string]bool{
  "electronics": true, "clothing": true, "home-appliances": true, "books": true, "furniture": true,
  "sports": true, "toys": true, "automotive": true, "new-arrival": true, "best-seller": true,
  "limited-edition": true, "exclusive": true, "eco-friendly": true, "handmade": true, "premium": true,
  "affordable": true, "summer-sale": true, "winter-collection": true, "black-friday": true,
  "holiday-special": true, "back-to-school": true, "sale": true, "clearance": true, "discounted": true,
  "buy-one-get-one": true, "free-shipping": true, "flash-deal": true, "new": true, "refurbished": true,
  "used": true, "vintage": true, "open-box": true, "kids": true, "adults": true, "women": true,
  "men": true, "seniors": true, "unisex": true, "nike": true, "samsung": true, "apple": true,
  "sony": true, "adidas": true, "puma": true, "lego": true, "durable": true, "lightweight": true,
  "waterproof": true, "ergonomic": true, "energy-efficient": true, "fast-charging": true,
  "multi-purpose": true, "high-performance": true, "red": true, "blue": true, "black": true,
  "white": true, "modern": true, "minimalistic": true, "classic": true,
}


func CreateUsersTable(db *sql.DB) {
  query := `
    CREATE TABLE IF NOT EXISTS users (
      id TEXT PRIMARY KEY,
      name TEXT,
      profileImage BLOB
      );
  `
  _, err := db.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Users Table created successfully!")
}

func CreateProductsTable(db *sql.DB) {
  query := `
    CREATE TABLE IF NOT EXISTS products (
      id TEXT PRIMARY KEY,
      name TEXT,
      price INTEGER,
      desc TEXT,
      quantity INTEGER,
      image BLOB
      );
  `
  _, err := db.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Product Table created successfully!")
}
func CreateTagsTable(db *sql.DB) {
  query := `
    CREATE TABLE IF NOT EXISTS tags (
      tagName TEXT PRIMARY KEY
      );
  `
  _, err := db.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Tags Table created successfully!")
}
func CreateTagsForProductTable(db *sql.DB) {
  query := `
    CREATE TABLE IF NOT EXISTS productTags (
      ProductId TEXT,
      TagName Test,
      FOREIGN KEY (ProductId) REFERENCES products(id),
      FOREIGN KEY (TagName) REFERENCES tags(tagName)
      UNIQUE (TagName, ProductId) 
      );
  `
  _, err := db.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Tags for Products Table created successfully!")
}

func InsertIntoUser(db *sql.DB, user wc.UserContext){
  checkQuery := `SELECT COUNT(*) FROM users WHERE id == ?`

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

  query := `INSERT INTO users (id, name, profileImage) VALUES (?, ?, ?)`

  _, err = db.Exec(query, user.UserID, user.UserName, user.ProfileImage)
  if err != nil {
      log.Fatal(err)
  }

  fmt.Println("User Data inserted successfully!")

}

func UpdateUserSes(db *sql.DB, sessionID string, userID string){
  query := `UPDATE sessions SET userID = ? WHERE id = ?`

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

func GetUser(db *sql.DB, sessionID string) wc.UserContext{
  query := `SELECT u.name, u.profileImage
  FROM sessions s
  JOIN users u ON s.UserID = u.id
  WHERE s.id = ?;`
  
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

func insertDefaultTags(db *sql.DB){
  query := `INSERT INTO tags (tagname) VALUES (?)`

  for _,tag := range(productTags){
    _, err := db.Exec(query, tag)
    if err != nil {
        fmt.Println("default Tags already created")
    }
  }

  fmt.Println("default Tags created")
}

func insertIntoProducts(db *sql.DB) {
  checkQuery := `SELECT COUNT(*) FROM products`

  row := db.QueryRow(checkQuery)
  var num int

  
  err := row.Scan(&num)
  if err != nil {
    log.Fatal(err)
  }

  if num >= 100 {
    println("Products are already inserted")
    return
  }
  
  query := `INSERT INTO products (id, name, price, desc, quantity, image) VALUES (?, ?, ?, ?, ?, ?)`

  imgData, err := os.ReadFile("test.png")
    if err != nil {
      log.Fatal(err)
    }
  
  for i := 0; i < 100; i++ {
    desc := "lorem Ipsum"
    uniqueID := uuid.New().String()
    if i > 50{ desc = "DEMO TEST"}
    _, err := db.Exec(query, uniqueID, "product " + strconv.Itoa(i), 25 + i, desc , 3, imgData)
    if err != nil {
        log.Fatal(err)
    }
  }

  fmt.Println("Data inserted successfully!")
}

func insertDefaultTags_for_Products(db *sql.DB){
	productsQuery := "SELECT id FROM products"
	rows, err := db.Query(productsQuery)
	if err != nil {
		fmt.Println("Error fetching products:", err)
		return
	}
	defer rows.Close()

	var products []string
	for rows.Next() {
		var productId string
		if err := rows.Scan(&productId); err != nil {
			fmt.Println("Error scanning product id:", err)
			return
		}
		products = append(products, productId)
	}

	tagsQuery := "SELECT tagName FROM tags"
	rows, err = db.Query(tagsQuery)
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		return
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			fmt.Println("Error scanning tag name:", err)
			return
		}
		tags = append(tags, tagName)
	}

  rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, productId := range products {
		// Select two random tags
		tag1 := tags[rand.Intn(len(tags))]
		tag2 := tags[rand.Intn(len(tags))]

		// Ensure both tags are different
		for tag1 == tag2 {
			tag2 = tags[rand.Intn(len(tags))]
		}

		// Step 4: Insert product-tag pairs into the productTags table
		insertQuery := `INSERT OR IGNORE INTO productTags (ProductId, TagName) VALUES (?, ?), (?, ?)`
		_, err := db.Exec(insertQuery, productId, tag1, productId, tag2)
		if err != nil {
			fmt.Println("Error inserting into productTags:", err)
			return
		}
		fmt.Printf("Assigned tags (%s, %s) to product %s\n", tag1, tag2, productId)
	}

	fmt.Println("Tag assignment complete!")

}

func GetProductsList(db *sql.DB,  offset int) ([]wc.Product, int) {
  query := `
    SELECT 
        COUNT(*) OVER() AS total,
        p.id, p.name, p.price, p.desc, p.quantity, p.image,
        GROUP_CONCAT(pt.tagName) AS tags
    FROM 
        products p
    LEFT JOIN 
        productTags pt ON p.id = pt.ProductId
    GROUP BY 
        p.id
    ORDER BY 
        p.price
    LIMIT 10
    OFFSET ?
`
      
  var ProductList []wc.Product

  rows, err := db.Query(query, offset)
  if err != nil {
    log.Fatal(err)
  }

  var id, name, desc, tagsStr string
  var price, total, quantity int
  var imgByte []byte

  for rows.Next() {
    err := rows.Scan(&total, &id, &name, &price, &desc, &quantity, &imgByte, &tagsStr)
    if err != nil {
        log.Fatal(err)
    }
    
    println("product list: adding:")
    println(id)
    println(name)
    imgStr := base64.StdEncoding.EncodeToString(imgByte)
    var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",") // Split by commas to get individual tags
		}
    ProductList = append(ProductList, wc.Product{
      Id: id, Name: name,
      Price: price, Desc: desc,
      Quantity: quantity, Image: imgStr,
      Tags: tags,
    })
  }
  println("product list total:")
  println(total)

  defer rows.Close()

  return ProductList, total

}



func GetProduct(db *sql.DB, prodId string) wc.Product {
  query := `SELECT id, name, price, desc, quantity, image FROM products WHERE id = ?`

  var product wc.Product

  rows, err := db.Query(query, prodId)
  if err != nil {
    log.Fatal(err)
  }

  var id, name, desc string
  var price, quantity int
  var imgByte []byte

  for rows.Next() {
    err := rows.Scan(&id, &name, &price, &desc, &quantity, &imgByte)
    if err != nil {
        log.Fatal(err)
    }
    
    imgStr := base64.StdEncoding.EncodeToString(imgByte)
    product = wc.Product{Id: id, Name: name, Price: price, Desc: desc, Quantity: quantity, Image: imgStr}
  }

    defer rows.Close()

  return product

}

func GetProductSearch(db *sql.DB, term string, offset int) ([]wc.Product, int) {
  terms := strings.Split(term, " ")

  //termsAmmount := len(terms)
  var query string
  var numOfTags = 0
  var numOfDesc = 0
  paramsForWHERE := make([]interface{}, 0, len(terms)*2)
  paramsForHAVING := make([]interface{}, 0, len(terms))
  paramsForHAVING_OR := make([]interface{}, 0, len(terms))

  for _, term := range terms{
    if productTagsMap[term] {
      numOfTags++
      term = "%" + term + "%"
      paramsForHAVING = append(paramsForHAVING, term)
    }else{
      numOfDesc++
      term = "%" + term + "%"
      paramsForWHERE = append(paramsForWHERE, term, term)
      paramsForHAVING_OR = append(paramsForHAVING_OR, term)
      // all the searched related to desc + not fully spelled tags
    }
  } 
  if ((numOfDesc >= 1) && (numOfTags >= 1)) {
    // TODO REPLACE UNION WITH LEFT/RIGHT UNION
    // SEARCH THOUGH TAGS, DESC, AND THINGS THAT WHERE INTENDED AS TAGS
    query = `
    WITH intersected AS (
      SELECT 
        p.id, 
        p.name, 
        p.price, 
        p.desc, 
        p.quantity, 
        p.image,
        GROUP_CONCAT(pt.tagName) AS tags
      FROM 
        products p
      LEFT JOIN 
        productTags pt ON p.id = pt.ProductId
      LEFT JOIN
        tags t ON pt.TagName = t.tagName
      WHERE (` + helperBuildWhereClause(numOfDesc) + `)
      GROUP BY 
        p.id
    
      INTERSECT
    
      SELECT 
        p.id, 
        p.name, 
        p.price, 
        p.desc, 
        p.quantity, 
        p.image,
        GROUP_CONCAT(pt.tagName) AS tags
      FROM 
        products p
      LEFT JOIN 
        productTags pt ON p.id = pt.ProductId
      LEFT JOIN
        tags t ON pt.TagName = t.tagName
      GROUP BY 
        p.id
      HAVING (` + helperBuildHavingClauseAND(numOfTags) + ` OR ` + helperBuildHavingClauseOR(numOfDesc) +  `)
    )
    SELECT 
      *,
      COUNT(*) OVER() AS total
    FROM intersected
    ORDER BY price 
    LIMIT 10
    OFFSET ?;
    `
  }else if numOfTags >= 1 {
    // SEARCH THOUGH TAGS
    query = `
    SELECT 
      p.id, p.name, p.price, p.desc, p.quantity, p.image,
      GROUP_CONCAT(pt.tagName) AS tags,
      COUNT(*) OVER() AS total
    FROM 
      products p
    LEFT JOIN 
      productTags pt ON p.id = pt.ProductId
    LEFT JOIN
      tags t ON pt.TagName = t.tagName
    GROUP BY 
      p.id
    HAVING (` + helperBuildHavingClauseAND(numOfTags) + `)
    ORDER BY price 
    LIMIT 10
    OFFSET ?
    `
  }else{
    // TODO REPLACE UNION WITH LEFT/RIGHT UNION
    // SEARCH THOUGH DESC AND CHECK IF USER MEANT TO TYPE A TAG
    query = `
    WITH intersected AS (
      SELECT 
        p.id, 
        p.name, 
        p.price, 
        p.desc, 
        p.quantity, 
        p.image,
        GROUP_CONCAT(pt.tagName) AS tags
      FROM 
        products p
      LEFT JOIN 
        productTags pt ON p.id = pt.ProductId
      LEFT JOIN
        tags t ON pt.TagName = t.tagName
      WHERE (` + helperBuildWhereClause(numOfDesc) + `)
      GROUP BY p.id
    
      UNION
    
      SELECT 
        p.id, 
        p.name, 
        p.price, 
        p.desc, 
        p.quantity, 
        p.image,
        GROUP_CONCAT(pt.tagName) AS tags
      FROM 
        products p
      LEFT JOIN 
        productTags pt ON p.id = pt.ProductId
      LEFT JOIN
        tags t ON pt.TagName = t.tagName
      GROUP BY p.id
      HAVING (` + helperBuildHavingClauseOR(numOfDesc) + ` )
    )
    SELECT 
      *,
      COUNT(*) OVER() AS total
    FROM intersected
    ORDER BY price 
    LIMIT 10
    OFFSET ?;
    `
  }
  

  //println("helper for WHERE")
  //println(helperBuildWhereClause(numOfDesc))
  //println("helper for HAVING")
  //println(helperBuildHavingClauseAND(numOfTags))
  //println(helperBuildHavingClauseOR(numOfDesc))
  //println("LEN IS ")
  //println(termsAmmount)
  //println("numOfDesc")
  //println(numOfDesc)
  //println("numOfTags")
  //println(numOfTags)
  //println("the final query")
  //fmt.Println(query)


  var params []interface{}
  params = append(paramsForWHERE, paramsForHAVING...)
  params = append(params, paramsForHAVING_OR...)

  params = append(params, offset)

  //println("THIS IS PARAMS")
  //
  //for _, param := range params {
  //    fmt.Println(param)
  //}

  var searchProductList []wc.Product

  rows, err := db.Query(query, params...)
  if err != nil {
    log.Fatal(err)
  }

  var id, name, desc, tagsStr string
  var total, price, quantity int
  var imgByte []byte

  for rows.Next() {
    err := rows.Scan(&id, &name, &price, &desc, &quantity, &imgByte, &tagsStr, &total)
    if err != nil {
        log.Fatal(err)
    }
    
    println("search: adding:")
    println(name)

    imgStr := base64.StdEncoding.EncodeToString(imgByte)
    var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",") // Split by commas to get individual tags
		}
    searchProductList = append(searchProductList, wc.Product{
      Id: id, Name: name,
      Price: price, Desc: desc,
      Quantity: quantity, Image: imgStr,
      Tags: tags,
    })
  }
  println("search total:")
  println(total)

  defer rows.Close()

  return searchProductList, total

}

func helperBuildWhereClause(termCount int) string{
  var clauses []string

  // Add conditions for 'name' and 'desc'
  for i := 0; i < termCount; i++ {
    clauses = append(clauses, "(p.name LIKE ? OR p.desc LIKE ?)")
  }

  return strings.Join(clauses, " OR ")

}
func helperBuildHavingClauseAND(termCount int) string {
  var clauses []string

  for i := 0; i < termCount; i++ {
    clauses = append(clauses, "(GROUP_CONCAT(t.tagName) LIKE ?)")
  }

  return strings.Join(clauses, " AND ")
}
func helperBuildHavingClauseOR(termCount int) string {
  var clauses []string

  for i := 0; i < termCount; i++ {
    clauses = append(clauses, "(GROUP_CONCAT(t.tagName) LIKE ?)")
  }
  return strings.Join(clauses, " OR ") 
}



type Session struct {
	ID                string
	UserID            string
	CreatedAt         int64
	CurrentPage       int64
  CurrentPageSearch int64
  Searching         bool
}

func CreateSessionsTable(db *sql.DB) {
  query := `
    CREATE TABLE IF NOT EXISTS sessions (
      id TEXT PRIMARY KEY,
      UserID,
      created_at INTEGER,
      current_page INTEGER,
      current_page_search INTEGER, 
      searching INTEGER,
      FOREIGN KEY (UserID) REFERENCES Users(id)
    );
  `
  _, err := db.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Session Table created successfully!")
}

func CreateSession(db *sql.DB) (string) {
	createdAt := time.Now().Unix()

  calc := createdAt + 223 + createdAt % 16

  sessionID := "se" + strconv.FormatInt(calc, 10) 

  query := `INSERT INTO sessions (id, UserID, created_at, current_page, current_page_search, searching) VALUES (?, ?, ?, ?, ?, ?)`
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


func IsLoggedIn(db *sql.DB, sessionID string) bool{
	query := `SELECT COUNT(UserID) FROM sessions WHERE id = ?`
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
func SelectCart_FirstTime(db *sql.DB, userloggedIn bool, sessionID string) []wc.CartItem {

  if userloggedIn == true {
    //[ ] check for the existing session(not used now) that had userID linked to it 
    query := `SELECT id 
      FROM sessions 
      WHERE UserID = (SELECT UserID FROM sessions WHERE id = ?) AND id != ? LIMIT 1`;

    var sessionID_OLD string
    var WeGotSessionThatHadThatAccLinked bool // really good naming
    // I really don't want to merge several sessionID's carts togerer (if there are some)
    row:= db.QueryRow(query, sessionID,sessionID)
    err := row.Scan(&sessionID_OLD)
    if err != nil {
      WeGotSessionThatHadThatAccLinked = false
    }else {
      WeGotSessionThatHadThatAccLinked = true
    }


    if WeGotSessionThatHadThatAccLinked {
      //[ ] copy that session's cart items to the new session

      // check if a new session started to have a cart, 
      //if so then we will ignore old session cart
      checkquery := `SELECT COUNT(*) 
        FROM cart 
        WHERE SessionId = ?;`

      var num int 
      row= db.QueryRow(checkquery, sessionID)
      err = row.Scan(&num)
      if err != nil {
        log.Fatal(err)
        return nil
      }
      if num == 0 {
        updatequery := `UPDATE cart
          SET SessionId = ?
          WHERE SessionId = ?;`

        _, err := db.Exec(updatequery, sessionID, sessionID_OLD)
        if err != nil {
          log.Fatal(err)
        }

      }
      //[ ] delete old session
      deletequery := ` DELETE FROM sessions 
        WHERE id = ?;`
      _, err = db.Exec(deletequery, sessionID_OLD)
      if err != nil {
        log.Fatal(err)
      }

    }
  }else{
    return SelectCart(db, sessionID)
  }
  return SelectCart(db, sessionID)
}


func SelectCart(db *sql.DB, sessionID string) []wc.CartItem {
  var rowData []wc.CartItem

  query := `SELECT 
      c.cartId,
      c.ProductId,
      p.name,
      p.price,
      c.quantity,
      p.desc,
      p.image,
      (p.price * c.quantity) AS total
      FROM cart c
      JOIN products p ON c.ProductId = p.id
      WHERE c.SessionId = ?`

  rows, err := db.Query(query, sessionID)
  if err != nil {
      log.Fatal(err)
      return nil
  }
  defer rows.Close()

  var cartId, productId, name, desc string
  var price, quantity, total int
  var imgByte []byte

  for rows.Next() {
    err := rows.Scan(&cartId, &productId, &name, &price, &quantity, &desc, &imgByte, &total)
    if err != nil {
      log.Fatal(err)
      return nil
    }

    imgStr := base64.StdEncoding.EncodeToString(imgByte)
    rowData = append(rowData, 
      wc.CartItem{
        Product: wc.Product{Id: productId, Name: name, Price: price, Quantity: quantity, Desc: desc, Image: imgStr}, 
        CartID: cartId,
        Total: total,
      })
  }

  if err := rows.Err(); err != nil {
    log.Fatal(err)
    return nil
  }

  return rowData
}
func CountFinalPrice(cartItemsList []wc.CartItem) (int) {

  var finalTotal int
  for _,item := range cartItemsList{
    finalTotal += item.Total
  }

  return finalTotal 
}



func SelectCartItem(db *sql.DB, productId string, sessionID string ) (wc.CartItem) {
    var product wc.CartItem

    query := `SELECT 
        c.cartId,
        p.name,
        p.price,
        c.quantity,
        (p.price * c.quantity) AS total
        FROM cart c
        JOIN products p ON c.ProductId = p.id
        WHERE c.ProductId = ? AND c.SessionID = ?`

    row := db.QueryRow(query, productId, sessionID)

    var cartId, name string
    var price, quantity, total int

    err := row.Scan(&cartId, &name, &price, &quantity, &total)
    if err != nil {
      log.Fatal(err)
      return product
    }

    product.Product = wc.Product{Id: productId, Name: name, Price: price, Quantity: quantity}
    product.CartID = cartId

    return product
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
func DeleteCart(db *sql.DB, sessionID string){

  query := `DELETE FROM cart WHERE SessionId = ?`

	_, err := db.Exec(query, sessionID)
  if err != nil {
    log.Fatal(err)
  }
}


func CreateDB() *sql.DB{

  db, err := sql.Open("sqlite3", "../database/products.db")
  if err != nil{
    log.Fatal(err)
  }

  CreateUsersTable(db)
  CreateProductsTable(db)
  CreateTagsTable(db)
  //insertDefaultTags(db)
  CreateTagsForProductTable(db)
  //insertDefaultTags_for_Products(db)
  CreateSessionsTable(db)
  CreateCartTable(db)
  insertIntoProducts(db)

  

  return db
}


