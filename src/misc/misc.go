package misc

import (
  wc "htmxNpython/web_context"
	"database/sql"
	db "htmxNpython/db_creation"
)

func NewPageContext(sqldb *sql.DB, values wc.InfiniteScroll, productList []wc.Product, isSearching bool, term string) wc.PageContext{
  return wc.PageContext{
    ProductsList: productList,
    Searching: isSearching,
    SearchTerm: term,
    Next: values.NewStart,
    More: values.More,
    NextProductsNums: values.NextProductsNums,
  }
}

func CreateCurentCart(sqldb *sql.DB, token string) wc.CurrentCart{

  items,_ := db.SelectCart(sqldb, token)
  println("CART CREATION!!!!!!!!!!")
  println(token)
  println(items)
  println("CART CREATION!!!!!!!!!!")

  return wc.CurrentCart{
    CartList: items,
  }
}
func CreateCurentCart_test(sqldb *sql.DB, token string) []wc.CartItem{

  items,_ := db.SelectCart(sqldb, token)
  println("CART CREATION!!!!!!!!!!")
  println(token)
  println(items)
  println("CART CREATION!!!!!!!!!!")

  return items
}



func NewSessionContext(db *sql.DB, token string, num int) wc.SessionContext{
  return wc.SessionContext{
    SessionID: token,
    CurrentPage: num,
    CurrentCart: CreateCurentCart(db, token),
  }
}

func NewProduct(id string, name string, price int, desc string, quantity int) wc.Product{
  return wc.Product{
    Id: id,
    Name: name,
    Price: price,
    Desc: desc,
    Quantity: quantity,
  }
}
func NewGlobalContext(sqldb *sql.DB, values wc.InfiniteScroll, token string, pagenum int, productsList []wc.Product, isSearching bool,term string) wc.GlobalContext{
  return wc.GlobalContext{
    PageContext: NewPageContext(sqldb, values, productsList, isSearching, term),
    SessionContext: NewSessionContext(sqldb, token, pagenum),
  } 
}
func NewGlobalContext_test(sqldb *sql.DB, session wc.Session_Test, page wc.PageContext_test, stripNums []int, productsList []wc.Product) wc.GlobalContext_Test{
  return wc.GenerateGlobalContext(session, page, productsList, stripNums, CreateCurentCart_test(sqldb,session.SessionID)) 
}

func GenerateNextProductNums(currentOffset int, itemsPerPage int, totalProducts int, searchTerm string) []int {
    totalPages := (totalProducts + itemsPerPage - 1) / itemsPerPage
    currentPage := currentOffset / itemsPerPage
    
    var values []int
    
    // Always show first page
    first_page := 0
    values = append(values, first_page)

    // Calculate window of pages around current
    startPage := max(1, currentPage-2)
    endPage := min(totalPages-1, currentPage+2)

    // Add pages in window
    for page := startPage; page <= endPage; page++ {
        values = append(values, page*itemsPerPage)
    }

    // Always show last page if not already included
    if endPage < totalPages-1 {
        values = append(values, (totalPages-1)*itemsPerPage)
    }

    return unique(values)
}


func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func unique(input []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    for _, val := range input {
        checkVal := val
        if !seen[checkVal] {
            seen[checkVal] = true
            result = append(result, val)
        }
    }
    return result
}


