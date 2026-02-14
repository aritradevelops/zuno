package handler

import "goserve/internal/pagination"

type Response[T any, I any, E any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
	Info    I      `json:"info"`
	Error   E      `json:"error"`
}
type NoInfo struct{}
type NoError struct{}

// @Description pagination options
type PaginationOptions struct {
	// the term to be searched for. default is ""
	Search string `query:"search" json:"search"`
	// whether to show deleted resources(only). defaults to false
	Trash bool `query:"trash" json:"trash"`
	// the field you want to sort the result set by. defaults to 'created_at'
	SortBy string `query:"sort_by" json:"sort_by"`
	// the sorting order. validation options are 'asc' | 'desc'. default to 'desc'
	SortOrder pagination.SortOrder `query:"sort_order" json:"sort_order"`
	// how many resource you want to fetch per set. default to 10
	Limit int `query:"limit" json:"limit"`
	// the page no you want fetch. default to 1.
	Page int `query:"page" json:"page"`
}

type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Info PaginationInfo `json:"info"`
}

// @Description pagination info
type PaginationInfo struct {
	// total resources count
	Total int `json:"total" example:"100"`
	// current page no
	Current int `json:"current" example:"1"`
	// whether next page exists
	HasNext bool `json:"has_next" example:"true"`
	// whether previous page exits
	HasPrev bool `json:"has_prev" example:"false"`
}

func Success[T any](message string, data T) Response[T, NoInfo, NoError] {
	return Response[T, NoInfo, NoError]{
		Message: message,
		Data:    data,
		Info:    NoInfo{},
		Error:   NoError{},
	}
}

func SuccessWithInfo[T any, I any](message string, data T, info I) Response[T, I, NoError] {
	return Response[T, I, NoError]{
		Message: message,
		Data:    data,
		Info:    info,
		Error:   NoError{},
	}
}

func Failure[E any](message string, err E) Response[struct{}, NoInfo, E] {
	return Response[struct{}, NoInfo, E]{
		Message: message,
		Data:    struct{}{},
		Info:    NoInfo{},
		Error:   err,
	}
}
