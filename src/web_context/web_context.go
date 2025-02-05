package web_context

import (
	"database/sql"

  m "htmxNpython/misc"
	db "htmxNpython/db_creation"
)

// \/\/ page content \/\/

// \/\/ PRODUCTS \/\/
// TODO STORE ID ON THE DB


func NewProduct(id int, name string, price int, desc string, quantity int) m.Product{
  return m.Product{
    Id: id,
    Name: name,
    Price: price,
    Desc: desc,
    Quantity: quantity,
  }
}

type ProductsList = []m.Product

func NewProductList(sqldb *sql.DB, from int, untill int) ProductsList{
  var products ProductsList
  productList := db.GetProductsList(sqldb, from, untill)
  diff := untill - from
  for i := 0; i < diff; i++ {

    products = append(products, NewProduct(productList[i].Id, productList[i].Name, productList[i].Price, productList[i].Desc, productList[i].Quantity))
  }
  return products
}


// \/\/ FOR THINGS THAT ARE SORED ON THE SERVER \/\/

type PageContext struct{
  ProductsList ProductsList
  Start int
  Next int
  More bool
  NextProductsNums []int
}

type InfiniteScroll struct {
  Start int
  NewStart int
  More bool
  NextProductsNums []int
}

func NewPageContext(sqldb *sql.DB, values InfiniteScroll) PageContext{
  return PageContext{
    ProductsList: NewProductList(sqldb, values.Start, values.NewStart),
    Start: values.Start,
    Next: values.NewStart,
    More: values.More,
    NextProductsNums: values.NextProductsNums,
  }
}

type CartItems struct{
  ProductsList ProductsList
  //Total int
} 

func NewCartItems(productsList ProductsList) CartItems{
  return CartItems{
    ProductsList: productsList,
    //Total: total,
  }
}

type CurrentCart struct{
  ProductsList ProductsList
}

func CreateCurentCart(sqldb *sql.DB, token string) CurrentCart{

  items,_ := db.SelectCart(sqldb, token)
  println("TEST id")

  return CurrentCart{
    ProductsList: items,
  }
}

type SessionContext struct{
  SessionID string
  CurrentPage int
  CurrentCart CurrentCart 
}

func NewSessionContext(db *sql.DB, token string, num int) SessionContext{
  return SessionContext{
    SessionID: token,
    CurrentPage: num,
    CurrentCart: CreateCurentCart(db, token),
  }
}

// \/\/ FOR THE WHOLE SITE \/\/
type GlobalContext struct {
  PageContext PageContext
  SessionContext SessionContext
}

func NewGlobalContext(sqldb *sql.DB, values InfiniteScroll, token string, pagenum int) GlobalContext{
  return GlobalContext{
    PageContext: NewPageContext(sqldb, values),
    SessionContext: NewSessionContext(sqldb, token, pagenum),
  } 
}

// ^^^^ page content ^^^^



