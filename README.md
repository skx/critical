[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/critical)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/critical)](https://goreportcard.com/report/github.com/skx/critical)


# criTiCaL

I had a bit of fun reading [TCL the Misunderstood](http://antirez.com/articoli/tclmisunderstood.html), and hacked up a simple TCL-like evaluator.


## Building/Using

Build in the usual way, for example:

```sh
$ go build .
```

Execute with the name of a TCL-file to execute:

```sh
    $ ./critical input.tcl
    A is set to: 4.000000
    Hello World
    4
    4
    ..
```


## Examples

The following is a simple example program which shows what the code here looks like:


```tcl
//
// Fibonacci sequence, written in the naive/recursive fashion.
//
proc fib {x} {
    if { expr $x <= 1 } {
        return 1
    } else {
        return [expr [fib [expr $x - 1]] + [fib [expr $x - 2]]]
    }
}


//
// Lets run this in a loop
//
set i 0
set max 20

while { expr $i <= $max } {
   puts "Fib $i is [fib $i]"
   incr i
}

```

Another example is the test-code which @antirez posted with his [picol writeup](http://oldblog.antirez.com/page/picol.html) which looks like this:

```tcl
proc square {x} {
    expr $x * $x
}

set a 1
while {expr $a <= 10} {
    if {expr $a == 5} {
        puts "\tMissing five!"
        incr a
        continue
        puts "After continue this won't be executed"
    }
    puts "I can compute that $a*$a = [square $a]"
    set a [expr $a + 1]
}

puts "I'm alive; after the 'while' loop."
```

This example is saved as [picol.tcl](picol.tcl) so you can run it from the repository:

```sh
go build . && ./critical ./picol.tcl
I can compute that 1*1 = 1
I can compute that 2*2 = 4
..
```


## Built In Commands

The following commands are available, and work as you'd expect:

* `break`, `continue`, `decr`, `expr`, `if`, `incr`, `proc`, `puts`, `return`, `set`, `while`.

In the near future we'll add `exit`, `for`, and a couple of the other simple primitives.

The complete list of standard [TCL commands](https://www.tcl.tk/man/tcl/TclCmd/contents.html) will almost certainly never be implemented, but pull-request to add omissions you need will be applied with thanks.



## Features

Read the file [input.tcl](input.tcl) to get a feel for the language, but in-brief you've got the following facilities available:

* Mathematical operations for `expr`
  * `+` `-` `/` `*` `<` `>` `<=` `>=`
  * Integer support only.  Sigh.
* Output to STDOUT via `puts`.
* Inline command expansion.
  * Including inside strings.
* Inline variable expansion.
  * Including inside strings.
* The ability to define procedures, via `proc`.

Badly implemented features:

* This fails as the spaces around `+` are necessary:
  * `puts [expr $a+$b]`
* Inline expansion swallows a character
  * `puts "[expr 3 + 3]ab` shows `3b` - where did `a` go?



## Testing

Most of the internal packages have high test-coverage, which can be exercised as you would expect:

```sh
$ go test ./...
```

In addition to the standard/static testing there are fuzz-based testers for the lexer and parser.  To run these run one of the following two sets of commands:

* **NOTE**: The fuzz-tester requires the use of golang version 1.18 or higher.


```sh
cd parser
go test -fuzztime=300s -parallel=1 -fuzz=FuzzParser -v
```

```sh
cd lexer
go test -fuzztime=300s -parallel=1 -fuzz=FuzzLexer -v

```
## Bugs?

Yeah I'm not surprised.  Please feel free to open a new issue **with your example** included so I can see how to fix it.


Steve
