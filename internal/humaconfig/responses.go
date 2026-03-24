package humaconfig

// PaginatedBody is the generic response body for paginated list endpoints.
type PaginatedBody[T any] struct {
	Items  []T   `json:"items"`
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}
