package repository

type ListOptions struct {
	Search    string
	Page      int
	PerPage   int
	SortBy    string
	SortOrder string
}

type ListInfo struct {
	Total   int
	Page    int
	PerPage int
}

type ListResponse[T any] struct {
	Data []T
	Info ListInfo
}
