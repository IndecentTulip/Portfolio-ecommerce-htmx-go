package misc

import (
	"database/sql"
	"fmt"
	db "HtmxReactGolang/db_api"
	wc "HtmxReactGolang/web_context"
	"strconv"
)

func NewPage(sqldb *sql.DB, session wc.SessionContext, page wc.PageContext) wc.GlobalContext{
  // PRODUCTSLIST
  var PRODUCTNUM int 
  var productsList []wc.Product
  const ITEMS_PER_PAGE = 20 // for both PRODUCTSLIST & NEXTPRODUCTSNUMS 

  if page.Is_Searching {
    productsList,PRODUCTNUM = db.GetProductListSearch(sqldb, page.SearchTerm, (page.Next - (ITEMS_PER_PAGE/2)))
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

  cart := db.SelectCart_FirstTime(sqldb, false, session.SessionID)

	userContext := wc.UserContext{
      UserName: "test",
      ProfileImage: "https://img.freepik.com/free-vector/blue-circle-with-white-user_78370-4707.jpg?t=st=1741218983~exp=1741222583~hmac=1b0ea872dd8d4b7b578200204a9df957dd072b79cd6b9644780d786ed6756b2b&w=740",
	}

  return wc.GenerateGlobalContext(
    productsList, nextProductsNums, cart,
    session, page, userContext) 
}



func NewGlobalContext(sqldb *sql.DB, session wc.SessionContext, page wc.PageContext) wc.GlobalContext{

  // USER
  userloggedIn := db.IsLoggedIn(sqldb, session.SessionID)

  var userContext wc.UserContext
  if !userloggedIn{
    userContext = wc.UserContext{
      UserName: "test",
      ProfileImage: "https://img.freepik.com/free-vector/blue-circle-with-white-user_78370-4707.jpg?t=st=1741218983~exp=1741222583~hmac=1b0ea872dd8d4b7b578200204a9df957dd072b79cd6b9644780d786ed6756b2b&w=740",
    }
  }else{
    userContext = db.GetUser(sqldb,session.SessionID) 
  } 
  // USER

  // PRODUCTSLIST
  var PRODUCTNUM int 
  var productsList []wc.Product
  const ITEMS_PER_PAGE = 20 // for both PRODUCTSLIST & NEXTPRODUCTSNUMS 

  if page.Is_Searching {
    productsList,PRODUCTNUM = db.GetProductListSearch(sqldb, page.SearchTerm, (page.Next - (ITEMS_PER_PAGE/2)))
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

  cart := db.SelectCart_FirstTime(sqldb, userloggedIn, session.SessionID)

  return wc.GenerateGlobalContext(
    productsList, nextProductsNums, cart,
    session, page, userContext) 
}

func BuildPaginationStrip(PRODUCTNUM, ITEMS_PER_PAGE, strip_num int) []wc.PaginationStrip {

		// TODO CONSIDER THAT pagination_strip items should be -1 to start from 0

		maxPagesToShow := 7

		// 100
		// 19
		// 20

		fmt.Println("😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥😥")
		totalPages := (PRODUCTNUM) / ITEMS_PER_PAGE 
		currentPage := strip_num + 1 // for UI
		var pagination_strip []wc.PaginationStrip
		fmt.Println(totalPages)
		fmt.Println(currentPage)

		// this works
		if totalPages <= maxPagesToShow {
      // Show all pages
      for i := 1; i <= totalPages; i++ {
          pagination_strip = append(pagination_strip, wc.PaginationStrip{Value: i-1, DisplayValue: strconv.Itoa(i)})
      }
			return pagination_strip
  	}

		// default first works
		pagination_strip = append(pagination_strip, wc.PaginationStrip{Value: 0, DisplayValue: "1"})
		
		//k
    if currentPage >= 4 {
      pagination_strip = append(pagination_strip,  wc.PaginationStrip{Value: -1, DisplayValue: "..."})
 // represents "..."
    }

    start := max(2, currentPage-1)
    end := min(totalPages-1, currentPage+1)
		// my addition for better beggining and end
		if currentPage == 1{
			end = start +1
		}
		if currentPage == totalPages{
			start = end - 1
		}

    for i := start; i <= end; i++ {
      pagination_strip = append(pagination_strip, wc.PaginationStrip{Value: i-1, DisplayValue: strconv.Itoa(i)})
    }

		// k
    if currentPage < totalPages-3 {
      pagination_strip = append(pagination_strip, wc.PaginationStrip{Value: -1, DisplayValue: "..."})
 // represents "..."
    }

		// default end
    pagination_strip = append(pagination_strip, wc.PaginationStrip{Value: totalPages-1, DisplayValue: strconv.Itoa(totalPages)})

    return pagination_strip
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


