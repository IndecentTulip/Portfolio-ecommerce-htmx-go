package main

import (
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

  items_per_page := 20
  productNum := 100

  e.GET("/", func(c echo.Context) error {
    // <><> session related <><>
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

    //<><>
    startStr := c.QueryParam("start")
    start, err := strconv.Atoi(startStr)
    currentPageNum,_ := db.GetSession(sqldb,sessionID)
    // 1 : 1 - 10
    // 2 : 20 - 30
    // 3 : 40 - 50
    page := int(int(currentPageNum.CurrentPage) - 1)
    range_start := page * items_per_page
    range_end := range_start + (items_per_page/2)

    // means we where not passed any start 
    // it will happen when:
    // 1. it's users first time
    // 2. they reloaded page(whitch may happen on diffent page ranges)
    if err != nil {
      start = range_start 
    }

    page_range := start >= range_start && start <= range_end  

    // <><>
    loadIndex := false
    if start == range_start {
      loadIndex = true
    }

    // <><>
    var newStart = start + (items_per_page/2)
    var more = newStart <= productNum
    var nextProductsNums []int
  

    // Used by the display of next pages 
    
    b := 1 
    t := 2
    // 60
    if start == range_start{
      if range_end + (items_per_page+(items_per_page/2)) >= productNum {
        b = 2
      }
      if start-items_per_page < 0{
        t = 3
      }

      if start > 1{
        for i,j := start-items_per_page, 0; i >= 0 && j < b; i, j = i-items_per_page, j+1{
          nextProductsNums = append([]int{i}, nextProductsNums...)
        }
      }
      nextProductsNums = append(nextProductsNums, start)
      for i,j := start+items_per_page, 0; i < productNum && j < t; i, j = i+items_per_page, j+1{
        nextProductsNums = append(nextProductsNums, i)

      }
    }
    println(sessionID)

    println(start)
    println(newStart)
    println(more)
    println(page_range)

    //db.GetProductsList(sqldb, start, newStart)

    if newStart > range_end{
      more = false
    } 
    var values wc.InfiniteScroll = wc.InfiniteScroll{
      Start: start,
      NewStart: newStart,
      More: more,
      NextProductsNums: nextProductsNums,
    }

    //sessionID="se1738281359"

    webContext := wc.NewGlobalContext(sqldb, values, sessionID, page+1)

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


  e.PUT("/tabnum", func(c echo.Context) error {

    numStr := c.FormValue("num")

    num,_ := strconv.Atoi(numStr)

    sessionID := c.Request().Header.Get("Cookie")
    if sessionID == ""{
      // TODO issue of not checking if token exists in db
      sessionID,_ = db.CreateSession(sqldb)
      db.UpdatePageNumSes(sqldb,sessionID,1)
    }else{
      sessionID = strings.Replace(sessionID, "session=", "",1)
    }

    currentPageNum,_ := db.GetSession(sqldb,sessionID)
    // 1 : 1 - 10
    // 2 : 20 - 30
    // 3 : 40 - 50
    page := int(int(currentPageNum.CurrentPage) - 1)
    range_start := page * items_per_page
    range_end := range_start + (items_per_page/2)

    println(sessionID)
    println(num)
    println(range_end)

    newPageNum := num / items_per_page
    newPageNum++
    
    db.UpdatePageNumSes(sqldb,sessionID,newPageNum)
    page = int(int(currentPageNum.CurrentPage) - 1)

    var sendContext any

    return c.Render(200, "restartpage", sendContext)
  });

  e.PUT("/addtocart", func(c echo.Context) error {

    productIDStr := c.FormValue("id")

    productID,_ := strconv.Atoi(productIDStr)

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


  e.Logger.Fatal(e.Start(":25258"))
}
