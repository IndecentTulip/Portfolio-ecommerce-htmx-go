package misc

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

type ProductNumsElement struct{   
  Num int
  SearchTerm string
}

func NewProduct(id string, name string, price int, desc string, quantity int) Product{
  return Product{
    Id: id,
    Name: name,
    Price: price,
    Desc: desc,
    Quantity: quantity,
  }
}

func GenerateNextProductNums(currentOffset int, itemsPerPage int, totalProducts int, searchTerm string) []ProductNumsElement {
    totalPages := (totalProducts + itemsPerPage - 1) / itemsPerPage
    currentPage := currentOffset / itemsPerPage
    
    var values []ProductNumsElement
    
    // Always show first page
    first_page := ProductNumsElement{Num: 0, SearchTerm: searchTerm}
    values = append(values, first_page)

    // Calculate window of pages around current
    startPage := max(1, currentPage-2)
    endPage := min(totalPages-1, currentPage+2)

    // Add pages in window
    for page := startPage; page <= endPage; page++ {
        temp := ProductNumsElement{Num: page*itemsPerPage, SearchTerm: searchTerm}
        values = append(values, temp)
    }

    // Always show last page if not already included
    if endPage < totalPages-1 {
        temp := ProductNumsElement{Num: (totalPages-1)*itemsPerPage, SearchTerm: searchTerm}
        values = append(values, temp)
    }

    return unique(values)
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func unique(input []ProductNumsElement) []ProductNumsElement {
    seen := make(map[int]bool)
    result := []ProductNumsElement{}
    for _, val := range input {
        checkVal := val.Num
        if !seen[checkVal] {
            seen[checkVal] = true
            result = append(result, val)
        }
    }
    return result
}


