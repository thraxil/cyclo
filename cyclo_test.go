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
	}
	for _, tc := range testCases {
		fset := token.NewFileSet()
		src, _ := ioutil.ReadFile(tc.File)
		f, _ := parser.ParseFile(fset, tc.File, src, 0)
		results := fileComplexity(f, fset)
		if results[0].Complexity != tc.Complexity {
			t.Error(fmt.Sprintf("wrong: %s %s", tc.File, tc.Function))
		}
	}
}
