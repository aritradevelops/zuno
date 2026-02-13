package service

import "fmt"

const ()

type Error interface {
	Long() string
	Short() string
	Code() int
	error
}

type serviceError struct {
	long  string
	short string
	code  int
}

func NewError(code int, long string, short string) Error {
	return &serviceError{
		long:  long,
		short: short,
		code:  code,
	}
}

// Code implements [Error].
func (s *serviceError) Code() int {
	return s.code
}

// Long implements [Error].
func (s *serviceError) Long() string {
	return s.long
}

// Short implements [Error].
func (s *serviceError) Short() string {
	return s.short
}

// Error implements [Error].
func (s *serviceError) Error() string {
	return fmt.Sprintf("Code: %d, Error: %s", s.code, s.long)
}
