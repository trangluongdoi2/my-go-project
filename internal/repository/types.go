package repository

type BaseQuery struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`

	SortBy    string `form:"sort_by" json:"sort_by"`       // e.g., "created_at", "name"
	SortOrder string `form:"sort_order" json:"sort_order"` // "asc" or "desc"

	Search string `form:"search" json:"search"` // General search across multiple fields
}

func (q BaseQuery) GetPage() int {
	if q.Page <= 0 {
		return 1
	}
	return q.Page
}

func (q BaseQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10
	}
	return q.Limit
}

func (q BaseQuery) GetSortFields() []string {
	if q.SortBy == "" {
		return []string{"created_at"}
	}
	return []string{q.SortBy}
}

func (q BaseQuery) GetSortOrders() []string {
	if q.SortOrder == "" {
		return []string{"desc"}
	}
	return []string{q.SortOrder}
}

type PreloadOption struct {
	Entity         string
	Conditions     map[string]interface{}
	SelectFields   []string
	Limit          int
	Offset         int
	NestedPreloads []PreloadOption
}

type Pagination[T any] struct {
	Total     int64 `json:"total"`
	Items     []T   `json:"items"`
	Limit     int   `json:"limit"`
	Page      int   `json:"page"`
	TotalPage int   `json:"total_page"`
}

type QueryOptions struct {
	SelectFields []string
	Limit        int
	Offset       int
	SortFields   []string
	SortOrders   []string
	Preloads     []PreloadOption
}
