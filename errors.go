package main

type ErrValidation string

func (e ErrValidation) Error() string {
	return "input validation failed"
}

type ErrDao string

func (e ErrDao) Error() string {
	return "dao has failed"
}
