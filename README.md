# cyclo

Simple cyclomatic complexity analysis for Go programs.

See: http://en.wikipedia.org/wiki/Cyclomatic_complexity

This program doesn't construct the full graph to calculate complexity;
it just does the simple count of the decision points.

## Installation

   go get github.com/thraxil/cyclo

(and make sure your `$GOPATH` is in your `$PATH`)

## Usage

    usage: cyclo [flags] [path ...]
      -max-complexity=0: max complexity

eg.

    $ cyclo cyclo.go
    cyclo.go:22:1:  main    1
    cyclo.go:27:1:  usage   1
    cyclo.go:33:1:  cycloMain       7
    cyclo.go:64:1:  getComplexity   2
    cyclo.go:70:1:  process 2
    cyclo.go:89:1:  processFile     11
    cyclo.go:130:1: visitFile       4
    cyclo.go:140:1: walkDir 1
    cyclo.go:144:1: report  1
    cyclo.go:149:1: isGoFile        2

    $ cyclo --max-complexity=10 cyclo.go
    cyclo.go:89:1:  processFile     11

## Bugs

I haven't yet figured out how to only count `return` statements that
aren't the last statement in the function, so functions with an
explicit `return` at the end currently return 1 higher than they
should.

I also have not really thought deeply yet about how `go` or `defer`
statements should be counted towards complexity. Anyone with good
ideas, let me know.

## License

BSD. 
