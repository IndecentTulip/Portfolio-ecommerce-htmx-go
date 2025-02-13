package misc

import (
  wc "htmxNpython/web_context"
	"database/sql"
	db "htmxNpython/db_creation"
)


func CreateCurentCart(sqldb *sql.DB, token string) []wc.CartItem{

  items,_ := db.SelectCart(sqldb, token)
  println("CART CREATION!!!!!!!!!!")
  println(token)
  println(items)
  println("CART CREATION!!!!!!!!!!")

  return items
}

func NewGlobalContext_test(sqldb *sql.DB, session wc.SessionContext, page wc.PageContext, stripNums []int, productsList []wc.Product) wc.GlobalContext{
  return wc.GenerateGlobalContext(session, page, productsList, stripNums, CreateCurentCart(sqldb,session.SessionID)) 
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


