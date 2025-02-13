package web_context

import "reflect"

type Product struct {
  Id string
  Name string
  Price int
  Desc string
  Quantity int
}

type CartItem struct{
  Product Product
  CartID  string
  Total   int
}

type PageContext struct{
  ProductsList []Product
  Next int
  More bool
  Searching bool
  SearchTerm string
  NextProductsNums []int
}

// used to set values
type InfiniteScroll struct {
  NewStart int
  More bool
  NextProductsNums []int
}

type CurrentCart struct{
  CartList []CartItem
}
type PageContext_test struct{
  Next int
  More bool
  SearchTerm string
}
type Session_Test struct {
  SessionID string
  CurrentPage int
}

type SessionContext struct{
  SessionID string
  CurrentPage int
  CurrentPageSearch int
  CurrentCart CurrentCart 
}

type GlobalContext struct {
  PageContext PageContext
  SessionContext SessionContext
}

type GlobalContext_Test struct {
	Values       map[string]interface{}
  ProductsList []Product
  NextProductsNums []int
  CartList      []CartItem
}

func GenerateGlobalContext(session Session_Test, page PageContext_test, productList []Product, strip []int, cartList []CartItem ) GlobalContext_Test {
  context := GlobalContext_Test{
    ProductsList: productList,
    NextProductsNums: strip,
    CartList: cartList,
    Values: make(map[string]interface{}),
  }

	sessionVal := reflect.ValueOf(session)

	for i := 0; i < sessionVal.NumField(); i++ {
		fieldName := sessionVal.Type().Field(i).Name
		fieldValue := sessionVal.Field(i).Interface()

    context.Values[fieldName] = fieldValue
	}
	pageVal := reflect.ValueOf(page)

	for i := 0; i < pageVal.NumField(); i++ {
		fieldName := pageVal.Type().Field(i).Name
		fieldValue := pageVal.Field(i).Interface()

    context.Values[fieldName] = fieldValue
	}

	return context
}


// GET RID OF THE HEADACKE, JUST MAKE EVER VAR GLOBAL
// AND ONLY USE ADDITONAL STRUCTS WHEN WORKING WITH ARRAY
// IT'S NOT FASTER THIS WAY, IT'S JUST MORE COMPLICATED FOR NO REASON
// OOP AND IT'S CONSEQUENCES FR FR

