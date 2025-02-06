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

func GetNextProductNums(sessionInfo db.Session, totalProducts int, isSearch bool) []int {
  const ITEMS_PER_PAGE = 20
  // MAKE IT NOT ONLY FOR SEARCH
  var currentOffset int
  if isSearch {
    currentOffset = (int(sessionInfo.CurrentPageSearch) - 1) * ITEMS_PER_PAGE
  }else{
    currentOffset = (int(sessionInfo.CurrentPage) - 1) * ITEMS_PER_PAGE
  }
  
  return m.GenerateNextProductNums(
      currentOffset,
      ITEMS_PER_PAGE,
      totalProducts,
  )
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
    
    var nextProductsNums []int 
    if !(start > range_start){
      nextProductsNums = GetNextProductNums(sessionInfo,PRODUCTNUM,false)
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

    webContext := wc.NewGlobalContext(sqldb, values, sessionID, page+1, productList)

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

    numStr := c.FormValue("num")

    num,_ := strconv.Atoi(numStr)

    sessionID := handleSessionWithoutAcc(sqldb,c)
    sessionInfo,_ := db.GetSession(sqldb,sessionID)
    page := int(int(sessionInfo.CurrentPage) - 1)
    range_start := page * ITEMS_PER_PAGE
    range_end := range_start + (ITEMS_PER_PAGE/2)

    println(sessionID)
    println(num)
    println(range_end)

    newPageNum := num / ITEMS_PER_PAGE
    newPageNum++
    
    db.UpdatePageNumSes(sqldb,sessionID,newPageNum)

    var sendContext any

    return c.Render(200, "restartpage", sendContext)
  });

  e.PUT("/addtocart", func(c echo.Context) error {

    productIDStr := c.FormValue("id")

    productID,_ := strconv.Atoi(productIDStr)

    // not using a function because in case someone not having sessionID
    // they should not be able to add anything to the cart in the first place
    sessionID := c.Request().Header.Get("Cookie")
    if sessionID == ""{
    }else{
      sessionID = strings.Replace(sessionID, "session=", "",1)
    }

    db.AddToCart(sqldb,sessionID,productID)

    productInfo,_ := db.SelectCartItem(sqldb,productID)

    type Product = m.Product

    type oobProduct struct{
      Product Product
      IsNew   bool
    }
    var isNew bool
    if productInfo.Quantity <= 1 {
      isNew = true 
    }else{
      isNew = false 
    }

    sendContext := oobProduct{
      Product: Product{
        Id: productInfo.Id,
        Name: productInfo.Name,
        Quantity: productInfo.Quantity,    
      },
      IsNew: isNew,
    } 

    return c.Render(200, "cartitems-oob", sendContext)
  });

  e.DELETE("/removefromcart", func(c echo.Context) error {
    cartIDStr := c.FormValue("id")

    cartID,_ := strconv.Atoi(cartIDStr)

    db.DeleteFromCart(sqldb, cartID)

    var sendContext any

    return c.Render(200, "temp", sendContext)
  });

  e.POST("/search", func(c echo.Context) error {
    // <><> session related <><>

    sessionID := handleSessionWithoutAcc(sqldb,c)

    // TODO CHANGE THIS
    startStr := c.QueryParam("start")
    start, err := strconv.Atoi(startStr)

    //<><>
    searchTerm := c.FormValue("search")
    searchTerm = strings.TrimSpace(searchTerm)  // Remove whitespace
    searchTerm = html.EscapeString(searchTerm) // Prevent XSS

    if searchTerm != ""{
      db.UpdateSearchingStatus(sqldb, sessionID, true)
    }else{
      db.UpdateSearchingStatus(sqldb,sessionID,false)
    }

    sessionInfo,_ := db.GetSession(sqldb,sessionID)
    productList,PRODUCTNUM := db.GetProductSearch(sqldb,searchTerm,10)

    nextProductsNums := GetNextProductNums(sessionInfo,PRODUCTNUM,true)

    println(PRODUCTNUM)

    page := int(int(sessionInfo.CurrentPageSearch) - 1)

    range_start := page * ITEMS_PER_PAGE
    range_end := range_start + (ITEMS_PER_PAGE/2)

    println("FOR SEARCH ranges")
    println(page)
    println(range_start)
    println(range_end)

    // TODO CHANGE THIS
    if err != nil {
      start = range_start 
    }

    page_range := start >= range_start && start <= range_end  

    println(page_range)

    // <><>
    loadIndex := false
    if start == range_start {
      loadIndex = true
    }

    // <><>
    var newStart = start + (ITEMS_PER_PAGE/2)
    var more = newStart <= PRODUCTNUM

    println("FOR SEARCH starts")
    println(start)
    println(newStart)

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

    webContext := wc.NewGlobalContext(sqldb, values, sessionID, page+1,productList)

    var sendContext any

    sendContext = webContext.PageContext

    template := "products"
    if loadIndex {
      template = "content"
      sendContext = webContext
    }

    if !page_range{
      template = "none"
    }

    return c.Render(200, template, sendContext)

  });

  e.Logger.Fatal(e.Start(":25258"))
}
