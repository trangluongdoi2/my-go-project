package repository

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
