package web_context

import "reflect"

type Product struct {
  Id string
  Name string
  Price int
  Desc string
  Quantity int
  Image string
  Tags []string
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
type UserContext struct {
  UserID string
  UserName string
  ProfileImage string
}

type GlobalContext struct {
	Values       map[string]interface{}
  ProductsList []Product
  NextProductsNums []int
  CartList      []CartItem
}

func GenerateGlobalContext(
  productList []Product, strip []int, cartList []CartItem,
  session SessionContext, page PageContext, user UserContext,
  //args ...interface{},
) GlobalContext {

  context := GlobalContext{
    ProductsList: productList,
    NextProductsNums: strip,
    CartList: cartList,
    Values: make(map[string]interface{}),
  }

  // TODO HARD CODE LOOPS INSDEAD OF sessionVal.NumField(); 
  //for _, arg := range args {
  //  val := reflect.ValueOf(arg)
  //
  //  for i := 0; i < val.NumField(); i++ {
  //    fieldName := val.Type().Field(i).Name
  //    fieldValue := val.Field(i).Interface()
  //
  //    context.Values[fieldName] = fieldValue
  //  }
  //}
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
	userVal := reflect.ValueOf(user)
  
	for i := 0; i < userVal.NumField(); i++ {
		fieldName := userVal.Type().Field(i).Name
		fieldValue := userVal.Field(i).Interface()
  
    context.Values[fieldName] = fieldValue
	}

	return context
}

