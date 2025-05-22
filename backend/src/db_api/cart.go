package db_api
import (
	"database/sql"
	"log"
	"encoding/base64"
	"github.com/google/uuid"

	wc "HtmxReactGolang/web_context"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"

)

type Cart struct{
  SessionId   string
  ProductId   string
}

func SelectCart_FirstTime(db *sql.DB, userloggedIn bool, sessionID string) []wc.CartItem {

  if userloggedIn == true {
    //[ ] check for the existing session(not used now) that had userID linked to it 
    query := `SELECT id 
      FROM sessions 
      WHERE UserID = (SELECT UserID FROM sessions WHERE id = $1) AND id != $1 LIMIT 1`;

    var sessionID_OLD string
    var WeGotSessionThatHadThatAccLinked bool // really good naming
    // I really don't want to merge several sessionID's carts togerer (if there are some)
    row:= db.QueryRow(query, sessionID)
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
        WHERE SessionId = $1;`

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
      WHERE c.SessionId = $1`

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
// TODO this should be in misc
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
        WHERE c.ProductId = $1 AND c.SessionID = $2`

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

  query := `INSERT INTO cart (CartId, SessionId, ProductId, Quantity)
		VALUES ($1, $2, $3, 1)
		ON CONFLICT (SessionId, ProductId) 
		DO UPDATE SET Quantity = cart.Quantity + 1;
	`

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

  query := `DELETE FROM cart WHERE CartId = $1`

	_, err := db.Exec(query, cartId)
  if err != nil {
    log.Fatal(err)
  }
}
func DeleteCart(db *sql.DB, sessionID string){

  query := `DELETE FROM cart WHERE SessionId = $1`

	_, err := db.Exec(query, sessionID)
  if err != nil {
    log.Fatal(err)
  }
}


