package web_context

import (
	"database/sql"

  m "htmxNpython/misc"
	db "htmxNpython/db_creation"
)

// \/\/ page content \/\/

// \/\/ PRODUCTS \/\/
// TODO STORE ID ON THE DB

type ProductsList = []m.Product


// \/\/ FOR THINGS THAT ARE SORED ON THE SERVER \/\/

type PageContext struct{
  ProductsList ProductsList
  Next int
  More bool
  Searching bool
  SearchTerm string
  NextProductsNums []m.ProductNumsElement
}

type InfiniteScroll struct {
  NewStart int
  More bool
  NextProductsNums []m.ProductNumsElement
}

func NewPageContext(sqldb *sql.DB, values InfiniteScroll, productList []m.Product, isSearching bool, term string) PageContext{
  return PageContext{
    ProductsList: productList,
    Searching: isSearching,
    SearchTerm: term,
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
  CartList []m.CartItem
}

func CreateCurentCart(sqldb *sql.DB, token string) CurrentCart{

  items,_ := db.SelectCart(sqldb, token)

  return CurrentCart{
    CartList: items,
  }
}

type SessionContext struct{
  SessionID string
  CurrentPage int
  CurrentPageSearch int
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

func NewGlobalContext(sqldb *sql.DB, values InfiniteScroll, token string, pagenum int, productsList []m.Product, isSearching bool,term string) GlobalContext{
  return GlobalContext{
    PageContext: NewPageContext(sqldb, values, productsList, isSearching, term),
    SessionContext: NewSessionContext(sqldb, token, pagenum),
  } 
}

// ^^^^ page content ^^^^


