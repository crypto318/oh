# A surprisingly powerful Unix shell

Oh is a Unix shell. If you've used other Unix shells, oh should feel
familiar.

![gif](img/oh.gif)

Where oh diverges from traditional Unix shells is in its programming
language features.

At its core, oh is a heavily modified dialect of the Scheme programming
language, complete with first-class continuations and proper tail
recursion. Like early Scheme implementations, oh exposes environments
as first-class values. Oh extends environments to allow both public and
private members and uses these extended first-class environments as the
basis for its prototype-based object system.

Written in Go, oh is also a concurrent programming language. It exposes
channels, in addition to pipes, as first-class values. As oh uses the
same syntax for code and data, channels and pipes can, in many cases, be
used interchangeably. This homoiconic nature also allows oh to support
fexprs which, in turn, allow oh to be easily extended. In fact, much of
oh is written in oh.<sup name="r1">[1](#f1)</sup>

For a detailed comparison to other Unix shells see: [Comparing oh to other Unix Shells](https://htmlpreview.github.io/?https://raw.githubusercontent.com/michaelmacinnis/oh/master/doc/comparison.html)

## Installing

With Go 1.5 or greater installed,

    go get github.com/michaelmacinnis/oh

According to [gox](https://github.com/mitchellh/gox), oh compiles on the
following platforms:

    darwin/386
    darwin/amd64
    dragonfly/amd64
    freebsd/386
    freebsd/amd64
    freebsd/arm
    linux/386
    linux/amd64
    linux/arm
    linux/arm64
    linux/ppc64
    linux/ppc64le
    netbsd/386
    netbsd/amd64
    netbsd/arm
    openbsd/386
    openbsd/amd64
    plan9/386
    plan9/amd64
    solaris/amd64
    windows/386
    windows/amd64

(Oh compiles and runs on Plan 9 and Windows but should be considered
experimental on those plarforms. On Solaris, interactive features are
limited).

## Using

For more detail see: [Using oh](doc/manual.md)

## License

[MIT](LICENSE)

<br><a name="f1">1</a>. [^](#r1) Currently, 
389 of 6699 lines of code are written in oh with an
additional 884 lines of Go generated by 175 lines
of oh.

