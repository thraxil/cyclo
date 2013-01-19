# cyclo

[![Build Status](https://travis-ci.org/thraxil/cyclo.png)](https://travis-ci.org/thraxil/cyclo)

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

If no flags are specified, it just finds every function definition in
the specified file(s) and displays its cyclomatic complexity:

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

The `--max-complexity` flag will filter the results so that only
functions with higher complexity are reported. It will also set the
error-code to a non-zero value if any functions trip it. This makes it
suitable for use in a git commit hook, for example.

    $ cyclo --max-complexity=10 cyclo.go
    cyclo.go:89:1:  processFile     11

The skeleton of the code was ripped off from `go fmt` (steal from the
best), so it should behave very similarly as far as how it globs
filenames, traverses directories, and so on.

## Bugs

I haven't yet figured out how to only count `return` statements that
aren't the last statement in the function, so functions with an
explicit `return` at the end currently return 1 higher than they
should.

## Notes

I have not really thought deeply yet about how `go` or `defer`
statements should be counted towards complexity. Anyone with good
ideas, let me know.

For now, I've made function literals count as a decision point towards
complexity. I'm not sure if that's right. Eg,

    myMap(mySlice, func (v int) {
        // do something with v
    })

My thinking is that often when you are using a function literal, it's
for that type of situation. Ie, "call this on some value later zero or
more times" and should be equivalent to a conditional or loop from a
complexity standpoint. Again, I'm open to arguments against this line
of reasoning.

The Python tool that inspired this has the ability to generate a graph
of the complexity map in dot format. That's cool, but I have never
found it useful for anything other than novelty, so I probably won't
bother implementing a similar feature here. If someone else is
interested in that, I'm happy to take a patch though.

## License

BSD.
