package main

type ErrValidation struct {
	message string
}

func NewErrValidation(message string) *ErrValidation {
	return &ErrValidation{
		message: message,
	}
}
func (e *ErrValidation) Error() string {
	return e.message
}

type ErrDao struct {
	message string
}

func NewErrDao(message string) *ErrDao {
	return &ErrDao{
		message: message,
	}
}
func (e *ErrDao) Error() string {
	return e.message
}
