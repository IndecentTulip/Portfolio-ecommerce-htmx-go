package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	db "htmxNpython/db_creation"
	tr "htmxNpython/temp_render"
	wc "htmxNpython/web_context"
	"io"
	"net/http"
	"os"

	m "htmxNpython/misc"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

func handleSessionWithoutAcc(sqldb *sql.DB, c echo.Context) string{
  sessionID := c.Request().Header.Get("Cookie")
  if sessionID == ""{
    sessionID = db.CreateSession(sqldb)
    db.UpdatePageNumSes(sqldb,sessionID,1)
  }else{
    sessionID = strings.Replace(sessionID, "session=", "",1)
    _,err := db.GetSession(sqldb,sessionID)
    if err != nil{
      sessionID = db.CreateSession(sqldb)
      fmt.Println(sessionID)
      db.UpdatePageNumSes(sqldb,sessionID,1)
    }
  }
  
  return sessionID
}

var oauth2Config_Google oauth2.Config
var oauth2Config_Github oauth2.Config
// made it random
var oauthStateString = "randomstate"

func init() {
	file, err := os.Open("authcred.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

  var clientIDGoogle, clientSecretGoogle string 
  var clientIDGithub, clientSecretGithub string 
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		clientIDGoogle = scanner.Text()
	} else {
		fmt.Println("No first line found.")
		return
	}

	if scanner.Scan() {
		clientSecretGoogle = scanner.Text()
	} else {
		fmt.Println("No second line found.")
		return
	}
	if scanner.Scan() {
		clientIDGithub = scanner.Text()
	} else {
		fmt.Println("No third line found.")
		return
	}
	if scanner.Scan() {
		clientSecretGithub = scanner.Text()
	} else {
		fmt.Println("No forth line found.")
		return
	}
  if scanner.Scan(){
    stripe.Key = scanner.Text()
  } else {
		fmt.Println("No fifth line found.")
		return
	}


	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading the file:", err)
	}

	oauth2Config_Google = oauth2.Config{
		ClientID:     clientIDGoogle,
		ClientSecret: clientSecretGoogle,
		RedirectURL:  "http://localhost:25258/callback",
		Scopes:       []string{"email", "profile"},   
		Endpoint:     google.Endpoint,               
	}
	oauth2Config_Github = oauth2.Config{
		ClientID:     clientIDGithub,
		ClientSecret: clientSecretGithub,
		RedirectURL:  "http://localhost:25258/callback/github",
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

}


func main(){

  sqldb := db.CreateDB()

  e := echo.New()
  e.Use(middleware.Logger())

  const ITEMS_PER_PAGE = 20

  // Renderer is an interface
  // by sayin that it is equal to "*Templates"
  // we would be able to use Render func that is made for that Struct
  e.Renderer = tr.NewTemplate()

  e.Static("/static/images", "images")
  e.Static("/static/css", "css")

  e.GET("/", func(c echo.Context) error {

    sessionID := handleSessionWithoutAcc(sqldb,c)
    
    startStr := c.QueryParam("start")
    start, err := strconv.Atoi(startStr)
    numStr := c.QueryParam("num")
    if numStr != ""{
      num,err := strconv.Atoi(numStr)
      if err != nil{
        num = 0
      }
      newPageNum := num / ITEMS_PER_PAGE
      newPageNum++
      db.UpdatePageNumSes(sqldb,sessionID,newPageNum)
    }

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

    page_range := start >= range_start && start <= range_end  

    var newStart = start + (ITEMS_PER_PAGE/2)

    _,PRODUCTNUM := db.GetProductsList(sqldb,start)

    var more = newStart <= PRODUCTNUM
  
    loadIndex := false
    if start == range_start {
      loadIndex = true
    }
    
    if newStart > range_end{
      more = false
    } 

    // TODO FIX

    ses := wc.SessionContext{
      SessionID: sessionID,
      CurrentPage: page+1,
    }
    pag := wc.PageContext{
      Next: newStart,
      More: more,
      Is_Searching: false,
      SearchTerm: "",
    }
    
    webContext := m.NewGlobalContext(sqldb,ses,pag)

    template := "products"
    if loadIndex {
      template = "index"
    }
    if !page_range{
      template = "none"
    }

    
    return c.Render(200, template, webContext)
  });

  e.GET("/search", func(c echo.Context) error {
    sessionID := handleSessionWithoutAcc(sqldb,c)

    startStr := c.QueryParam("start")
    start, err := strconv.Atoi(startStr)
    numStr := c.QueryParam("num")
    if numStr != ""{
      num,err := strconv.Atoi(numStr)
      if err != nil{
        num = 0
      }
      newPageNum := num / ITEMS_PER_PAGE
      newPageNum++
      db.UpdatePageSearchNumSes(sqldb,sessionID,newPageNum)
    }

    // maybe it would be better idea to expect newSearch to be ether 0 or 1
    // rather then relying on the error to happen...
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

    page := int(int(sessionInfo.CurrentPageSearch) - 1)

    range_start := page * ITEMS_PER_PAGE
    range_end := range_start + (ITEMS_PER_PAGE/2)

    if err != nil {
      start = range_start 
    }

    _,PRODUCTNUM := db.GetProductSearch(sqldb,searchTerm,start)

    page_range := start >= range_start && start <= range_end  

    fmt.Println(page_range)

    var newStart = start + (ITEMS_PER_PAGE/2)
    var more = newStart <= PRODUCTNUM

    loadIndex := false
    // newSearch should be 0, if it is true
    // idk I found my own implementation confusing
    if (start == range_start) && newSearch != 1 {
      loadIndex = true
    }

    if newStart > range_end{
      more = false
    } 

    ses := wc.SessionContext{
      SessionID: sessionID,
      CurrentPage: page+1,
    }
    pag := wc.PageContext{
      Next: newStart,
      More: more,
      Is_Searching: true,
      SearchTerm: searchTerm,
    }

    webContext := m.NewGlobalContext(sqldb,ses,pag)

    template := "products"
    if loadIndex {
      template = "index"
    }

    if !page_range{
      template = "none"
    }

    return c.Render(200, template, webContext)

  });

  e.PUT("/addtocart", func(c echo.Context) error {
    // WE CAN GET THE SESSION 
    // AND THE ITEM ID FROM THIS REQUEST
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
    fmt.Println("!!!! API /addtocart !!!!")
    println(sessionID)

    db.AddToCart(sqldb,sessionID,productID)

    // FIX
    productInfo := db.SelectCartItem(sqldb,productID,sessionID)

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

  // !!!!!!!!!!!!!!!!!!!!!!!!
  e.GET("/p/:product_id", func(c echo.Context) error {

    productId := c.Param("product_id")
    sessionID := c.Request().Header.Get("Cookie")
    sessionID = strings.Replace(sessionID, "session=", "",1)
    fmt.Println("VISITING PRODUCT")
    fmt.Println(productId)
    fmt.Println(sessionID)

    type WebContext struct{
      Product wc.Product
      Values struct{
        SessionID string
      }
      CartList []wc.CartItem
    }

    product :=  db.GetProduct(sqldb,productId)
    //session := wc.NewSessionContext(sqldb, sessionID, 0)

    cartList := db.SelectCart(sqldb,sessionID)
    webContext := WebContext{
      Product: product,
      Values: struct{SessionID string}{
        SessionID: sessionID,
      },
      CartList: cartList,
    }

    return c.Render(200,"productPage", webContext)
  });

  e.GET("/c/:session_id", func(c echo.Context) error {

    sessionID := c.Param("session_id")
    fmt.Println("BYING")
    fmt.Println(sessionID)

    cartInfo := db.SelectCart(sqldb,sessionID)
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

  e.GET("/payment", func(c echo.Context) error {

    sessionID := c.Request().Header.Get("Cookie")
    if sessionID == ""{
    }else{
      sessionID = strings.Replace(sessionID, "session=", "",1)
    }

    cartInfo := db.SelectCart(sqldb,sessionID)
    finalPrice := db.CountFinalPrice(cartInfo)

    // TODO store finalPrice as a float number
    fmt.Println(finalPrice)

    params := &stripe.CheckoutSessionParams{
      PaymentMethodTypes: stripe.StringSlice([]string{
        "card",
      }),
      LineItems: []*stripe.CheckoutSessionLineItemParams{
        {
          Name:        stripe.String("Product Name"),
          Description: stripe.String("Description of the product"),
          Amount:      stripe.Int64(int64(finalPrice * 100)),
          Currency:    stripe.String(string(stripe.CurrencyCAD)),
          Quantity:    stripe.Int64(1),
        },
      },
      Mode:      stripe.String(string(stripe.CheckoutSessionModePayment)),
      SuccessURL: stripe.String("http://localhost:25258/success"),
      CancelURL:  stripe.String("http://localhost:25258/"),
    }

    session, err := session.New(params)
    if err != nil {
      return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create session: %v", err))
    }

    return c.Redirect(http.StatusSeeOther, session.URL)

  });

  e.GET("/login", func(c echo.Context) error {

    typeOfAuth := c.QueryParam("type")

    var url string
    if typeOfAuth == "google"{
      url = oauth2Config_Google.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
    }else if typeOfAuth == "github"{
      url = oauth2Config_Github.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
    }

    //http.Redirect(w, r, oauth2Config.AuthCodeURL(oauth2State, oauth2.AccessTypeOffline), http.StatusFound)
    return c.Redirect(http.StatusFound, url)
  });

  // will leave it as just callback
  // because it's OG and I want to shoot myself in the foot in the future
  e.GET("/callback", func(c echo.Context) error {

    state := c.QueryParam("state")

    fmt.Println("state")
    fmt.Println(state)

    code := c.QueryParam("code")

    println("CODE CODE CODE CODE CODE")
    println(code)
    if code == "" {
      return c.String(http.StatusBadRequest, "Code not found")
    }
    token, err := oauth2Config_Google.Exchange(c.Request().Context(), code)
    if err != nil {
      return c.String(http.StatusInternalServerError, "Error during Exchange with oauth2")
    }
    client := oauth2Config_Google.Client(c.Request().Context(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
    if err != nil {
      return c.String(http.StatusInternalServerError, "Error getting userinfo")
    }

    //decoder := json.NewDecoder(resp.Body)
    //var userInfo UserInfo
    var userInfo map[string]interface{}
    body, err := io.ReadAll(resp.Body)

    err = json.Unmarshal(body, &userInfo)

    sessionID := c.Request().Header.Get("Cookie")
    if sessionID == ""{
    }else{
      sessionID = strings.Replace(sessionID, "session=", "",1)
    }

    userContext := wc.UserContext{
      UserName: userInfo["name"].(string),
      UserID: userInfo["sub"].(string),
      ProfileImage: userInfo["picture"].(string),
    }

    db.InsertIntoUser(sqldb, userContext)
    db.UpdateUserSes(sqldb,sessionID,userContext.UserID)
    println("CALLED UpdateUserSes")

    return c.Redirect(http.StatusSeeOther, "http://localhost:25258/")
  });

  e.GET("/callback/github", func(c echo.Context) error {

    state := c.QueryParam("state")
    if state != oauthStateString {
      return c.String(http.StatusBadRequest, "Invalid state")
    }

    code := c.QueryParam("code")
    if code == "" {
      return c.String(http.StatusBadRequest, "Code not found")
    }

    token, err := oauth2Config_Github.Exchange(c.Request().Context(), code)
    if err != nil {
      return c.String(http.StatusInternalServerError, "Failed to exchange token")
    }

    // issue was here
    client := github.NewClient(oauth2.NewClient(c.Request().Context(), oauth2.StaticTokenSource(&oauth2.Token{
      AccessToken: token.AccessToken,
    })))

    user, _, err := client.Users.Get(c.Request().Context(), "")
    if err != nil {
      return c.String(http.StatusInternalServerError, "Failed to get user info")
    }

    sessionID := c.Request().Header.Get("Cookie")
    if sessionID == ""{
    }else{
      sessionID = strings.Replace(sessionID, "session=", "",1)
    }

    userContext := wc.UserContext{
      UserName: *user.Login,
      UserID: strconv.FormatInt(user.GetID(), 10),
      ProfileImage: *user.AvatarURL,
    }

    db.InsertIntoUser(sqldb, userContext)
    db.UpdateUserSes(sqldb,sessionID,userContext.UserID)
    println("CALLED UpdateUserSes")

    return c.Redirect(http.StatusSeeOther, "http://localhost:25258/")

    });

  e.Logger.Fatal(e.Start(":25258"))
}
