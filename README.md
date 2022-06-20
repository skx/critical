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

Badly implemented features:

* Inline command expansion only works to one level
  * This is fine: `puts [expr 3 + 3]`
  * This fails: `puts [expr 3 + [expr 1 + 2]]`
* This fails as the spaces around `+` are necessary:
  * `puts [expr $a+$b]`
* Inline expansion swallows a character
  * `puts "[expr 3 + 3]ab` shows `3b` - where did `a` go?

Missing features?

* The ability to define `procs`.


## Bugs?

Yeah I'm not surprised.  Please feel free to open a new issue **with your example** included so I can see how to fix it.


Steve
