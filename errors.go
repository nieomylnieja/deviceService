package main

import "log"

type ErrValidation string

func (e ErrValidation) Error() string {
	return "input validation failed"
}

type ErrDao string

func (e ErrDao) Error() string {
	return "dao has failed"
}

func panicOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s :%s", msg, err)
	}
}
