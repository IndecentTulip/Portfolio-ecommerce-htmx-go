package misc

import (
  wc "htmxNpython/web_context"
	"database/sql"
	db "htmxNpython/db_creation"
)

func NewGlobalContext(sqldb *sql.DB, session wc.SessionContext, page wc.PageContext, user wc.UserContext) wc.GlobalContext{
  const ITEMS_PER_PAGE = 20 
  var productsList []wc.Product
  var PRODUCTNUM int 

  if page.Is_Searching {
    productsList,PRODUCTNUM = db.GetProductSearch(sqldb, page.SearchTerm, (page.Next - (ITEMS_PER_PAGE/2)))
  }else{
    productsList,PRODUCTNUM = db.GetProductsList(sqldb,(page.Next - (ITEMS_PER_PAGE/2)))
  }

  pageNum := session.CurrentPage - 1
  range_start := pageNum * ITEMS_PER_PAGE

  var nextProductsNums []int
  if !((page.Next - 10) > range_start){
    searchTerm := ""
    nextProductsNums = GetNextProductNums(session,PRODUCTNUM,false,searchTerm)
  }
 
  return wc.GenerateGlobalContext(
    productsList, nextProductsNums, CreateCurentCart(sqldb,session.SessionID),
    session, page, user) 
}

func CreateCurentCart(sqldb *sql.DB, token string) []wc.CartItem{

  items,_ := db.SelectCart(sqldb, token)
  //println("CART CREATION!!!!!!!!!!")
  //println(token)
  //println(items)
  //println("CART CREATION!!!!!!!!!!")

  return items
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


