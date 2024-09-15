// Use of this source code is governed by a BSD-style license

package main // import "github.com/thraxil/cyclo"

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	exitCode      = 0
	maxComplexity = flag.Int("max-complexity", 0, "max complexity")
)

func main() {
	cycloMain()
	os.Exit(exitCode)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: cyclo [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func cycloMain() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		if err := processFile("<standard input>",
			os.Stdin, os.Stdout, true); err != nil {
			report(err)
		}
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, nil, os.Stdout, false); err != nil {
				report(err)
			}
		}
	}
}

type fcomplexity struct {
	complexity int
}

func (f fcomplexity) getComplexity() int {
	return f.complexity
}

// quick and dirty count of if's, for's, case's, etc.
// not accurate, but already useful
func (f *fcomplexity) process(x ast.Node) {
	ast.Inspect(x, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.BranchStmt:
			f.complexity++
		case *ast.SwitchStmt:
			f.complexity++
		case *ast.FuncLit:
			f.complexity++
		case *ast.ForStmt:
			f.complexity++
		case *ast.RangeStmt:
			f.complexity++
		case *ast.IfStmt:
			f.complexity++
		case *ast.ReturnStmt:
			// how to only count this if it's not at the end of the function?
			f.complexity++
		}
		return true
	})
}

type complexityResult struct {
	Position     token.Position
	FunctionName *ast.Ident
	Complexity   int
}

func fileComplexity(f ast.Node, fset *token.FileSet) []complexityResult {
	results := []complexityResult{}
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			fc := fcomplexity{complexity: 1}
			fc.process(x)
			complexity := fc.getComplexity()
			results = append(results, complexityResult{
				Position:     fset.Position(n.Pos()),
				FunctionName: x.Name,
				Complexity:   complexity,
			})
		}
		return true
	})
	return results
}

func processFile(filename string, in io.Reader,
	out io.Writer, stdin bool) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		panic(err)
	}

	results := fileComplexity(f, fset)
	for _, r := range results {
		if r.Complexity > *maxComplexity {
			fmt.Printf("%s:\t%s\t%d\n", r.Position, r.FunctionName, r.Complexity)
			if *maxComplexity != 0 {
				exitCode = 1
			}
		}
	}
	return err
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGoFile(f) {
		err = processFile(path, nil, os.Stdout, false)
	}
	if err != nil {
		report(err)
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}
