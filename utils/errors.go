package utils

import (
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func CheckErrorWithData[T any](data T, err error) T {
	CheckError(err)
	return data
}
