// internal/subscription/parser/utils.go
package parser

import "strconv"

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
