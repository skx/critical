[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/critical)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/critical)](https://goreportcard.com/report/github.com/skx/critical)
[![license](https://img.shields.io/github/license/skx/critical.svg)](https://github.com/skx/critical/blob/master/LICENSE)


# criTiCaL

I had a bit of fun reading [TCL the Misunderstood](http://antirez.com/articoli/tclmisunderstood.html), and hacked up a simple TCL-like evaluator.

The process was described a little in my blog here:

* [Writing a simple TCL interpreter in golang](https://blog.steve.fi/writing_a_simple_tcl_interpreter_in_golang.html)



## Building & Usage

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

The interpreter contains an embedded "standard-library", which you can view at [stdlib/stdlib.tcl](stdlib/stdlib.tcl), which is loaded along with any file that you specify.

To disable the use of the standard library run:

```sh
   $ ./critical -no-stdlib path/to/file.tcl
```



## Examples

The following is a simple example program which shows what the code here looks like:


```tcl
//
// Fibonacci sequence, written in the naive/recursive fashion.
//
proc fib {x} {
    if { <= $x 1 } {
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

while {<= $i $max } {
   puts "Fib $i is [fib $i]"
   incr i
}

```

Another example is the test-code which @antirez posted with his [picol writeup](http://oldblog.antirez.com/page/picol.html) which looks like this:

```tcl
proc square {x} {
    * $x $x
}

set a 1
while {<= $a 10} {
    if {== $a 5} {
        puts {Missing five!}
        set a [+ $a 1]
        continue
    }
    puts "I can compute that $a*$a = [square $a]"
    set a [+ $a 1]
}
```

This example is contained within this repository as [picol.tcl](picol.tcl), so you can run it directly:

```sh
go build . && ./critical ./picol.tcl
I can compute that 1*1 = 1
I can compute that 2*2 = 4
..
```


## Built In Commands

The following commands are available, and work as you'd expect:

* `break`, `continue`, `decr`, `exit`, `expr`, `if`, `incr`, `proc`, `puts`, `return`, `set`, `while`.

In the near future we'll add `for`, and a couple of the other simple primitives.

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



## Testing

Most of the internal packages have high test-coverage, which can be exercised as you would expect:

```sh
$ go test ./...
```

In addition to the standard/static testing there are fuzz-based testers for the lexer and parser.  To run these run one of the following two sets of commands:

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
