package utils

import (
	"log"
)

func NewCheck(name string) *checker {
	return &checker{name: name}
}

type checker struct {
	name string
}

func CheckError(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func CheckErrorWithData[T any](data T, err error) T {
	CheckError(err)
	return data
}
