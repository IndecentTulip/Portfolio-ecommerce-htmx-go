package main

import (
	"html/template"
	"io"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// \/\/ echo setup \/\/

type Templates struct{
  templates *template.Template
}

func newTemplate() *Templates {
  return &Templates{
    templates: template.Must(template.ParseGlob("views/*.html")),
  }
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data) 
}

// ^^^^ echo setup ^^^^

// \/\/ page content \/\/

// \/\/ PRODUCTS \/\/
// TODO STORE ID ON THE DB
type Product struct {
  Id int
  Name string
  Prise int
  Desc string
}

func newProduct(id int, name string, prise int, desc string) Product{
  return Product{
    Id: id,
    Name: name,
    Prise: prise,
    Desc: desc,
  }
}

type ProductsList = []Product

func newProductList(from int, untill int) ProductsList{
  var products ProductsList
  for i := from; i < untill; i++ {
    products = append(products, newProduct(i, "product " + strconv.Itoa(i), 10 + i, "lorem ipsum"))
  }
  return products
}

// \/\/ FOR THINGS THAT ARE SORED ON THE SERVER \/\/

type PageContext struct{
  ProductsList ProductsList
  Start int
  Next int
  More bool
}

type InfiniteScroll struct {
  start int
  newStart int
  more bool
}

func newPageContext(values InfiniteScroll) PageContext{
  return PageContext{
    ProductsList: newProductList(values.start, values.newStart),
    Start: values.start,
    Next: values.newStart,
    More: values.more,
  }
}

// \/\/ FOR THE WHOLE SITE \/\/
type GlobalContext struct {
  PageContext PageContext
}

func newGlobalContext(values InfiniteScroll) GlobalContext{
  return GlobalContext{
    PageContext: newPageContext(values),
  } 
}

// ^^^^ page content ^^^^


func main(){

  e := echo.New()
  e.Use(middleware.Logger())

  e.Renderer = newTemplate()

  e.Static("/static/images", "images")
  e.Static("/static/css", "css")

  // TODO
  // read about OOB
  e.GET("/", func(c echo.Context) error {
    startStr := c.QueryParam("start")
    start, err := strconv.Atoi(startStr)
    if err != nil {
        start = 0
    }

    var newStart = start + 10
    var more = newStart < 100 

    //println(start)
    //println(newStart)
    //println(more)

    var values InfiniteScroll = InfiniteScroll{
      start: start,
      newStart: newStart,
      more: more,
    }

    Context := newGlobalContext(values)

    var sendContext any

    sendContext = Context.PageContext
    template := "products"
    if start == 0 {
      template = "index"
      sendContext = Context
    }

    return c.Render(200, template, sendContext)
  });


  e.Logger.Fatal(e.Start(":25252"))
}
