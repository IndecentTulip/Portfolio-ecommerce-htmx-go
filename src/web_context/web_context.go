package web_context

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

type Session struct {
	ID                string
	UserID            int
	CreatedAt         int64
	CurrentPage       int64
  CurrentPageSearch int64
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

// GET RID OF THE HEADACKE, JUST MAKE EVER VAR GLOBAL
// AND ONLY USE ADDITONAL STRUCTS WHEN WORKING WITH ARRAY
// IT'S NOT FASTER THIS WAY, IT'S JUST MORE COMPLICATED FOR NO REASON
// OOP AND IT'S CONSEQUENCES FR FR

