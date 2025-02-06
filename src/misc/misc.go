package misc

type Product struct {
  Id int
  Name string
  Price int
  Desc string
  Quantity int
}
// TODO ADD QUANTITY

type Session struct {
	ID                string
	UserID            int
	CreatedAt         int64
	CurrentPage       int64
  CurrentPageSearch int64
}

func NewProduct(id int, name string, price int, desc string, quantity int) Product{
  return Product{
    Id: id,
    Name: name,
    Price: price,
    Desc: desc,
    Quantity: quantity,
  }
}

func GenerateNextProductNums(currentOffset int, itemsPerPage int, totalProducts int) []int {
    totalPages := (totalProducts + itemsPerPage - 1) / itemsPerPage
    currentPage := currentOffset / itemsPerPage
    var nums []int

    // Always show first page
    nums = append(nums, 0)

    // Calculate window of pages around current
    startPage := max(1, currentPage-2)
    endPage := min(totalPages-1, currentPage+2)

    // Add pages in window
    for page := startPage; page <= endPage; page++ {
        nums = append(nums, page*itemsPerPage)
    }

    // Always show last page if not already included
    if endPage < totalPages-1 {
        nums = append(nums, (totalPages-1)*itemsPerPage)
    }

    return unique(nums)
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

func unique(input []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    for _, val := range input {
        if !seen[val] {
            seen[val] = true
            result = append(result, val)
        }
    }
    return result
}


