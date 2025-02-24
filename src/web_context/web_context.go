package web_context

import "reflect"

type Product struct {
  Id string
  Name string
  Price int
  Desc string
  Quantity int
  Image string
}
type CartItem struct{
  Product Product
  CartID  string
  Total   int
}

type PageContext struct{
  Next int
  More bool
  SearchTerm string
  Is_Searching bool
}

type SessionContext struct {
  SessionID string
  CurrentPage int
}

type GlobalContext struct {
	Values       map[string]interface{}
  ProductsList []Product
  NextProductsNums []int
  CartList      []CartItem
}

func GenerateGlobalContext(session SessionContext, page PageContext, productList []Product, strip []int, cartList []CartItem ) GlobalContext {
  context := GlobalContext{
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

