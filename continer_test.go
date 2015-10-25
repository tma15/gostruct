package gostruct

import (
	"fmt"
	"testing"
)

func TestMatrix(t *testing.T) {
	m := NewMatrix()
	fmt.Println(m)
	m.Resize(1, 0) /* 2 rows, 1 colum */

	fmt.Println(m)
	fmt.Println(len(m), len(m[0]))

	m[1][0] += 1.
	fmt.Println(m)
}
