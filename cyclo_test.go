package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"
)

type testCase struct {
	File string
	Function string
	Complexity int
}

func TestProcessFile(t *testing.T) {
	testCases := []testCase{
		{"test_cases/1.go", "a", 1},
		{"test_cases/2.go", "b", 2},
		{"test_cases/3.go", "a", 2},
		{"test_cases/4.go", "a", 3},
	}
	for _, tc := range testCases {
		fset := token.NewFileSet()
		src, _ := ioutil.ReadFile(tc.File)
		f, _ := parser.ParseFile(fset, tc.File, src, 0)
		results := fileComplexity(f, fset)
		if results[0].Complexity != tc.Complexity {
			t.Error(fmt.Sprintf("wrong: %s %s [%d vs %d]", tc.File, tc.Function,
				tc.Complexity, results[0].Complexity))
		}
	}
}
