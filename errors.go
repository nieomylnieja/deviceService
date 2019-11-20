package main

type ErrValidation string

func (e ErrValidation) Error() string {
	return "input validation failed"
}

type ErrDao string

func (e ErrDao) Error() string {
	return "dao has failed"
}

type ErrNotFound string

func (e ErrNotFound) Error() string {
	return "device doesn't exist"
}
