package common

import "fmt"

var (
	currentId = 0
)

func UniqueIdIdentifier() string {
	c := currentId
	currentId++
	return fmt.Sprintf("__________%x", c)
}
