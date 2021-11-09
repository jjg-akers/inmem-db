package test

import (
	"fmt"
	"testing"
)

// Coverage reports if required coverage was met
func Coverage(code int, goal float64, required bool) int {
	if testing.CoverMode() == "" {
		return code
	}

	c := testing.Coverage()

	if c < goal {
		fmt.Printf("coverage: %.1f%% below goal %.1f%%\n", c*100.0, goal*100.0)

		if required {
			fmt.Println("FAIL\tdue to insufficient test coverage")
			code = 1
		}
	}

	return code
}
