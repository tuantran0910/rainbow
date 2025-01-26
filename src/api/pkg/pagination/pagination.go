package pagination

import "math"

type Pagination struct {
	Limit       int `json:"limit"`
	CurrentPage int `json:"current_page"`
	TotalItems  int `json:"total_items"`
	LastPage    int `json:"last_page"`
	Offset      int `json:"offset"`
}

func NewPagination(page, limit, totalItems int) *Pagination {
	// Calculate the total pages and offset
	totalPages := getTotalPages(totalItems, limit)
	offset := getOffset(page, limit)

	return &Pagination{
		Limit:       limit,
		CurrentPage: page,
		TotalItems:  totalItems,
		LastPage:    totalPages,
		Offset:      offset,
	}
}

func getTotalPages(totalItems, limit int) int {
	return int(math.Ceil(float64(totalItems) / float64(limit)))
}

func getOffset(page, limit int) int {
	return (page - 1) * limit
}
