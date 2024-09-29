package script

import (
)

type Script struct {
	Tokens []*Token
	Warnings []string

	StartAddress int
	StackAddress int

	Labels map[int]string
}
