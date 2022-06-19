# criTiCaL

I had a bit of fun reading [TCL the Misunderstood](http://antirez.com/articoli/tclmisunderstood.html), and hacked up a simple TCL-like evaluator.


## Building/Using

Build in the usual way, for example:

```sh
$ go build .
```

Execute with zero arguments to run the embedded examples, otherwise
pass the name of a file:

```sh
    $ ./critical input.tcl
    A is set to: 4.000000
    Hello World
    4
    4
    ..
```


## Features

Allows only basic operations:

* Mathematical operations for `expr`
  * `+` `-` `/` `*` `<` `>` `<=` `>=`
  * Integer support only.  Sigh.
* Simple output via `puts`
* Inline command expansion.
* Inline variable expansion.

Badly implemented features:

* if
* while
* block support `{` `}`

Missing features:

* No real parser
