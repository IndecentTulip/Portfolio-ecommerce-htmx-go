package main

import (
	"database/sql"
	"fmt"
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

func GetNextProductNums(sessionInfo db.Session, totalProducts int, isSearch bool, searchTerm string) []int {
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
    fmt.Println("FOR / ranges")
    fmt.Println(page)
    fmt.Println(range_start)
    fmt.Println(range_end)

    page_range := start >= range_start && start <= range_end  

    fmt.Println(page_range)

    // <><>
    var newStart = start + (ITEMS_PER_PAGE/2)

    productList,PRODUCTNUM := db.GetProductsList(sqldb,start)

    var more = newStart <= PRODUCTNUM
  
    fmt.Println("FOR / starts")
    fmt.Println(start)
    fmt.Println(newStart)
    
    var nextProductsNums []int
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

    fmt.Println(sessionID)
    fmt.Println(page_range)

    fmt.Println(more)

    ses := wc.Session_Test{
      SessionID: sessionID,
      CurrentPage: page+1,
    }
    pag := wc.PageContext_test{
      Next: newStart,
      More: more,
      Is_Searching: false,
      SearchTerm: "",
    }
    webContext := m.NewGlobalContext_test(sqldb,ses,pag,nextProductsNums,productList)

    template := "products"
    if loadIndex {
      template = "index"
    }

    if !page_range{
      template = "none"
    }

    return c.Render(200, template, webContext)
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

    fmt.Println(sessionID)
    fmt.Println(num)
    fmt.Println(searchTerm)
    fmt.Println(range_end)

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
    fmt.Println("!!!! API /addtocart !!!!")
    fmt.Println(productID)
    // not using a function because in case someone not having sessionID
    // they should not be able to add anything to the cart in the first place
    sessionID := c.Request().Header.Get("Cookie")
    if sessionID == ""{
    }else{
      sessionID = strings.Replace(sessionID, "session=", "",1)
    }

    db.AddToCart(sqldb,sessionID,productID)

    productInfo,_ := db.SelectCartItem(sqldb,productID)

    type oobProduct struct{
      CartItem wc.CartItem
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

  e.PUT("/addtocart/page", func(c echo.Context) error {

    productID := c.QueryParam("id")
    fmt.Println("!!!! API /addtocart !!!!")
    fmt.Println(productID)
    sessionID := c.Request().Header.Get("Cookie")
    if sessionID == ""{
    }else{
      sessionID = strings.Replace(sessionID, "session=", "",1)
    }

    db.AddToCart(sqldb,sessionID,productID)

    var sendContext any

    return c.Render(200, "temp", sendContext)
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
    fmt.Println(searchTerm)
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
    fmt.Println(PRODUCTNUM)



    var nextProductsNums []int 
    if !(start > range_start){
      nextProductsNums = GetNextProductNums(sessionInfo,PRODUCTNUM,true, searchTerm)
    }

    fmt.Println("FOR SEARCH ranges")
    fmt.Println(page)
    fmt.Println(range_start)
    fmt.Println(range_end)

    // TODO CHANGE THIS

    page_range := start >= range_start && start <= range_end  

    fmt.Println(page_range)

    // <><>
    // <><>
    var newStart = start + (ITEMS_PER_PAGE/2)
    var more = newStart <= PRODUCTNUM

    fmt.Println("FOR SEARCH starts")
    fmt.Println(start)
    fmt.Println(newStart)

    loadIndex := false
    if (start == range_start) && newSearch != 1 {
      loadIndex = true
    }

    if newStart > range_end{
      more = false
    } 

    fmt.Println(sessionID)
    fmt.Println(page_range)

    fmt.Println(more)

    ses := wc.Session_Test{
      SessionID: sessionID,
      CurrentPage: page+1,
    }
    pag := wc.PageContext_test{
      Next: newStart,
      More: more,
      Is_Searching: true,
      SearchTerm: searchTerm,
    }
    webContext := m.NewGlobalContext_test(sqldb,ses,pag,nextProductsNums,productList)


    template := "products"
    if loadIndex {
      template = "index"
    }

    if !page_range{
      template = "none"
    }

    return c.Render(200, template, webContext)

  });

  e.GET("/p/:product_id", func(c echo.Context) error {

    productId := c.Param("product_id")
    sessionID := c.Request().Header.Get("Cookie")
    sessionID = strings.Replace(sessionID, "session=", "",1)
    fmt.Println("VISITING PRODUCT")
    fmt.Println(productId)
    fmt.Println(sessionID)

    type WebContext struct{
      Product wc.Product
      SessionContext wc.SessionContext
    }
    fmt.Println("TEST")
    product :=  db.GetProduct(sqldb,productId)
    //session := wc.NewSessionContext(sqldb, sessionID, 0)

    webContext := WebContext{
      Product: product,
 //     SessionContext: session,
    }

    return c.Render(200,"productPage", webContext)
  });

  e.GET("/c/:session_id", func(c echo.Context) error {

    sessionID := c.Param("session_id")
    fmt.Println("BYING")
    fmt.Println(sessionID)

    cartInfo,_ := db.SelectCart(sqldb,sessionID)
    finalPrice := db.CountFinalPrice(cartInfo)

    type CartPage struct{
      CartInfo []wc.CartItem
      FinalPrice int
    }

    cartContext := CartPage{
      CartInfo: cartInfo,
      FinalPrice: finalPrice,
    } 

    return c.Render(200,"cartPage", cartContext)
  });

  e.PUT("/payment", func(c echo.Context) error {

    cardnumber := c.FormValue("cardNumber")

    fmt.Println("!!!!!!!!!!!!!!!!!!!")
    fmt.Println(cardnumber)
    fmt.Println("!!!!!!!!!!!!!!!!!!!")

    //stripe.Key = "sk_test_4eC39HqLyjWDarjtT1zdp7dc"
    //var paymentData struct {
    //    Amount int64  `json:"amount"`
    //    Token  string `json:"token"` // Token from the frontend
    //}
    //err := json.NewDecoder(r.Body).Decode(&paymentData)
    //if err != nil {
    //    http.Error(w, "Invalid data", http.StatusBadRequest)
    //    return
    //}
    //
    //// Create a payment intent with Stripe
    //pi, err := paymentintent.New(&stripe.PaymentIntentParams{
    //    Amount:   stripe.Int64(paymentData.Amount),
    //    Currency: stripe.String(string(stripe.CurrencyUSD)),
    //    PaymentMethod: stripe.String(paymentData.Token),
    //    ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodManual)),
    //    Confirm: stripe.Bool(true),
    //})
    //
    //if err != nil {
    //    http.Error(w, "Payment processing error", http.StatusInternalServerError)
    //    return
    //}
    //
    //// Return the payment intent result
    //response := map[string]interface{}{
    //    "success": pi.Status == stripe.PaymentIntentStatusSucceeded,
    //}

    var sendContext any

    return c.Render(200,"restartpage", sendContext)
  });




  e.Logger.Fatal(e.Start(":25258"))
}
