package misc

import (
  wc "htmxNpython/web_context"
	"database/sql"
	db "htmxNpython/db_creation"
)

func NewGlobalContext(sqldb *sql.DB, session wc.SessionContext, page wc.PageContext) wc.GlobalContext{

  // USER
  println("TEST 1")
  userloggedIn := db.IsLoggedIn(sqldb, session.SessionID)

  var userContext wc.UserContext
  if !userloggedIn{
    userContext = wc.UserContext{
      UserName: "test",
      ProfileImage: "https://img.freepik.com/free-vector/blue-circle-with-white-user_78370-4707.jpg?t=st=1741218983~exp=1741222583~hmac=1b0ea872dd8d4b7b578200204a9df957dd072b79cd6b9644780d786ed6756b2b&w=740",
    }
  }else{
    println("TEST 2")
    userContext = db.GetUser(sqldb,session.SessionID) 
    println("TEST 3")
  } 
  // USER

  // PRODUCTSLIST
  var PRODUCTNUM int 
  var productsList []wc.Product
  const ITEMS_PER_PAGE = 20 // for both PRODUCTSLIST & NEXTPRODUCTSNUMS 

  if page.Is_Searching {
    println("TEST 4")
    productsList,PRODUCTNUM = db.GetProductSearch(sqldb, page.SearchTerm, (page.Next - (ITEMS_PER_PAGE/2)))
    println("TEST 5")
  }else{
    productsList,PRODUCTNUM = db.GetProductsList(sqldb,(page.Next - (ITEMS_PER_PAGE/2)))
  }
  // PRODUCTSLIST

  pageNum := session.CurrentPage - 1
  range_start := pageNum * ITEMS_PER_PAGE

  // NEXTPRODUCTSNUMS
  var nextProductsNums []int
  if !((page.Next - 10) > range_start){
    searchTerm := ""
    nextProductsNums = GetNextProductNums(session,PRODUCTNUM,false,searchTerm)
  }
  // NEXTPRODUCTSNUMS

  println("TEST 6")
  cart := db.SelectCart_FirstTime(sqldb, userloggedIn, session.SessionID)
  println("TEST 7")

  return wc.GenerateGlobalContext(
    productsList, nextProductsNums, cart,
    session, page, userContext) 
}

func GetNextProductNums(sessionInfo wc.SessionContext, totalProducts int, isSearch bool, searchTerm string) []int {
  const ITEMS_PER_PAGE = 20
  // MAKE IT NOT ONLY FOR SEARCH
  var currentOffset int
  currentOffset = (int(sessionInfo.CurrentPage) - 1) * ITEMS_PER_PAGE
  
  return GenerateNextProductNums(currentOffset, ITEMS_PER_PAGE, totalProducts, searchTerm)
}

func GenerateNextProductNums(currentOffset int, itemsPerPage int, totalProducts int, searchTerm string) []int {
    totalPages := (totalProducts + itemsPerPage - 1) / itemsPerPage
    currentPage := currentOffset / itemsPerPage
    
    var values []int
    
    first_page := 0
    values = append(values, first_page)

    startPage := max(1, currentPage-2)
    endPage := min(totalPages-1, currentPage+2)

    for page := startPage; page <= endPage; page++ {
        values = append(values, page*itemsPerPage)
    }

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


