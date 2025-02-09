package main

import (
	"database/sql"
	"html"
	db "htmxNpython/db_creation"
	tr "htmxNpython/temp_render"
	wc "htmxNpython/web_context"

	m "htmxNpython/misc"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func GetNextProductNums(sessionInfo db.Session, totalProducts int, isSearch bool, searchTerm string) []m.ProductNumsElement {
  const ITEMS_PER_PAGE = 20
  // MAKE IT NOT ONLY FOR SEARCH
  var currentOffset int
  if isSearch {
    currentOffset = (int(sessionInfo.CurrentPageSearch) - 1) * ITEMS_PER_PAGE
  }else{
    currentOffset = (int(sessionInfo.CurrentPage) - 1) * ITEMS_PER_PAGE
  }
  
  return m.GenerateNextProductNums(currentOffset, ITEMS_PER_PAGE, totalProducts, searchTerm)
}

func handleSessionWithoutAcc(sqldb *sql.DB, c echo.Context) string{
  sessionID := c.Request().Header.Get("Cookie")
  if sessionID == ""{
    sessionID,_ = db.CreateSession(sqldb)
    db.UpdatePageNumSes(sqldb,sessionID,1)
  }else{
    sessionID = strings.Replace(sessionID, "session=", "",1)
    _,err := db.GetSession(sqldb,sessionID)
    if err != nil{
      sessionID,_ = db.CreateSession(sqldb)
      db.UpdatePageNumSes(sqldb,sessionID,1)
    }
  }
  
  return sessionID
}

func main(){

  sqldb := db.CreateDB()

  e := echo.New()
  e.Use(middleware.Logger())

  // Renderer is an interface
  // by sayin that it is equal to "*Templates"
  // we would be able to use Render func that is made for that Struct
  e.Renderer = tr.NewTemplate()

  e.Static("/static/images", "images")
  e.Static("/static/css", "css")

  ITEMS_PER_PAGE := 20

  e.GET("/", func(c echo.Context) error {
    // <><> session related <><>
    sessionID := handleSessionWithoutAcc(sqldb,c)
    //<><>
    startStr := c.QueryParam("start")
    start, err := strconv.Atoi(startStr)
    sessionInfo,_ := db.GetSession(sqldb,sessionID)

    // 1 : 1 - 10
    // 2 : 20 - 30
    // 3 : 40 - 50
    page := int(int(sessionInfo.CurrentPage) - 1)
    db.UpdatePageSearchNumSes(sqldb,sessionID,1)
    db.UpdateSearchingStatus(sqldb,sessionID,false)
    range_start := page * ITEMS_PER_PAGE
    range_end := range_start + (ITEMS_PER_PAGE/2)

    // means we where not passed any start 
    // it will happen when:
    // 1. it's users first time
    // 2. they reloaded page(whitch may happen on diffent page ranges)
    if err != nil {
      start = range_start 
    }
    println("FOR / ranges")
    println(page)
    println(range_start)
    println(range_end)

    page_range := start >= range_start && start <= range_end  

    println(page_range)

    // <><>
    var newStart = start + (ITEMS_PER_PAGE/2)

    productList,PRODUCTNUM := db.GetProductsList(sqldb,start)

    var more = newStart <= PRODUCTNUM
  
    println("FOR / starts")
    println(start)
    println(newStart)
    
    var nextProductsNums []m.ProductNumsElement 
    if !(start > range_start){
      searchTerm := ""
      nextProductsNums = GetNextProductNums(sessionInfo,PRODUCTNUM,false,searchTerm)
    }

    // <><><>
    loadIndex := false
    if start == range_start {
      loadIndex = true
    }
    
    if newStart > range_end{
      more = false
    } 
    // <><><>

    println(sessionID)
    println(page_range)

    println(more)

    var values wc.InfiniteScroll = wc.InfiniteScroll{
      NewStart: newStart,
      More: more,
      NextProductsNums: nextProductsNums,
    }

    webContext := wc.NewGlobalContext(sqldb, values, sessionID, page+1, productList, false,"")

    var sendContext interface{}

    sendContext = webContext.PageContext

    template := "products"
    if loadIndex {
      template = "index"
      sendContext = webContext
    }

    if !page_range{
      template = "none"
    }

    return c.Render(200, template, sendContext)
  });


  e.PUT("/tabnum", func(c echo.Context) error {

    numStr := c.QueryParam("num")
    searchTerm := c.QueryParam("search")

    num,err := strconv.Atoi(numStr)
    if err != nil{
      num = 0
    }

    sessionID := handleSessionWithoutAcc(sqldb,c)
    sessionInfo,_ := db.GetSession(sqldb,sessionID)
    page := int(int(sessionInfo.CurrentPage) - 1)
    range_start := page * ITEMS_PER_PAGE
    range_end := range_start + (ITEMS_PER_PAGE/2)

    println(sessionID)
    println(num)
    println(searchTerm)
    println(range_end)

    newPageNum := num / ITEMS_PER_PAGE
    newPageNum++
    
    type SendContext struct{
      IsSearching bool
      SearchTerm string
    }
    isSearching := false
    if searchTerm != ""{
      db.UpdatePageSearchNumSes(sqldb,sessionID,newPageNum)
      isSearching = true
    }else{
      db.UpdatePageNumSes(sqldb,sessionID,newPageNum)
    }

    sendContext := SendContext{
      IsSearching: isSearching,
      SearchTerm: searchTerm,
    }

    return c.Render(200, "restartpage", sendContext)
  });

  e.PUT("/addtocart", func(c echo.Context) error {

    productID := c.QueryParam("id")
    println("!!!! API /addtocart !!!!")
    println(productID)
    // not using a function because in case someone not having sessionID
    // they should not be able to add anything to the cart in the first place
    sessionID := c.Request().Header.Get("Cookie")
    if sessionID == ""{
    }else{
      sessionID = strings.Replace(sessionID, "session=", "",1)
    }

    db.AddToCart(sqldb,sessionID,productID)

    productInfo,_ := db.SelectCartItem(sqldb,productID)

    type CartItem = m.CartItem

    type oobProduct struct{
      CartItem CartItem
      IsNew   bool
    }
    var isNew bool
    if productInfo.Product.Quantity <= 1 {
      isNew = true 
    }else{
      isNew = false 
    }

    sendContext := oobProduct{
      CartItem: productInfo,
      IsNew: isNew,
    } 

    return c.Render(200, "cartitems-oob", sendContext)
  });

  e.DELETE("/removefromcart", func(c echo.Context) error {
    productID := c.QueryParam("id")

    db.DeleteFromCart(sqldb, productID)

    var sendContext any

    return c.Render(200, "temp", sendContext)
  });

  e.GET("/search", func(c echo.Context) error {
    sessionID := handleSessionWithoutAcc(sqldb,c)

    startStr := c.QueryParam("start")
    start, err := strconv.Atoi(startStr)

    newSearchStr := c.QueryParam("newSearch")
    newSearch,erro := strconv.Atoi(newSearchStr)
    if erro != nil{
      newSearch = 0
    }

    //<><>
    searchTerm := c.QueryParam("search")
    println(searchTerm)
    searchTerm = strings.TrimSpace(searchTerm)  // Remove whitespace
    searchTerm = html.EscapeString(searchTerm) // Prevent XSS

    if searchTerm != ""{
      db.UpdateSearchingStatus(sqldb, sessionID, true)
    }else{
      db.UpdateSearchingStatus(sqldb,sessionID,false)
    }

    sessionInfo,_ := db.GetSession(sqldb,sessionID)
    // page is needed to understand if you need to lazy load more content
    page := int(int(sessionInfo.CurrentPageSearch) - 1)

    range_start := page * ITEMS_PER_PAGE
    range_end := range_start + (ITEMS_PER_PAGE/2)

    // start is needed to understand what is the current OFFSET for the SELECT query is 
    if err != nil {
      start = range_start 
    }

    productList,PRODUCTNUM := db.GetProductSearch(sqldb,searchTerm,start)
    println(PRODUCTNUM)



    var nextProductsNums []m.ProductNumsElement 
    if !(start > range_start){
      nextProductsNums = GetNextProductNums(sessionInfo,PRODUCTNUM,false, searchTerm)
    }

    println("FOR SEARCH ranges")
    println(page)
    println(range_start)
    println(range_end)

    // TODO CHANGE THIS

    page_range := start >= range_start && start <= range_end  

    println(page_range)

    // <><>
    // <><>
    var newStart = start + (ITEMS_PER_PAGE/2)
    var more = newStart <= PRODUCTNUM

    println("FOR SEARCH starts")
    println(start)
    println(newStart)

    loadIndex := false
    if (start == range_start) && newSearch != 1 {
      loadIndex = true
    }

    if newStart > range_end{
      more = false
    } 

    println(sessionID)
    println(page_range)

    println(more)

    var values wc.InfiniteScroll = wc.InfiniteScroll{
      NewStart: newStart,
      More: more,
      NextProductsNums: nextProductsNums,
    }

    webContext := wc.NewGlobalContext(sqldb, values, sessionID, page+1,productList, true, searchTerm)

    var sendContext any

    sendContext = webContext.PageContext

    template := "products"
    if loadIndex {
      template = "index"
      sendContext = webContext
    }

    if !page_range{
      template = "none"
    }

    return c.Render(200, template, sendContext)

  });

  e.GET("/c/:cart_id", func(c echo.Context) error {

    cartID := c.Param("cart_id")
    println(cartID)

    var sendContext any

    return c.Render(404,"index", sendContext)
  });



  e.Logger.Fatal(e.Start(":25258"))
}
