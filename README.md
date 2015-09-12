# go-where -- Find where a Go identifier is defined

Library and a command-line utility to locate Go definitions.

```console
$ go-where github.com/tv42/where#Ident
/home/tv/go/src/github.com/tv42/where/where.go:30:6
```

So far, the Go community has not standardized on a way to type a full
import and identifiers into a single string. Using dot `.` as a
separator between the import path and identifiers (as opposed to
package *names* and identifiers) is ambiguous;
`example.com/foo.bar.baz`, we can't tell where the repository part
ends.

To avoid ambiguity, we choose a separator character that is not found
in import paths: the hash sign, `#`. In practice, you can think of the
input as "what [godoc.org](http://godoc.org/) would take".

The input syntax is `IMPORTPATH[#IDENTIFIER[.FIELD]...[.METHOD]]`, or
simplifying `IMPORTPATH[#NAME[.NAME]...]`.

`go-where` is primarily intended for uses like embedding stable links
to source into other files -- see below for Emacs `org-mode`
integration.

As such, more value is placed on clarity over conciseness, and we do
not ambiguous input; `go doc exec.Cmd.Run` may happen to find the
right thing for simple cases, but `go-where os/exec#Cmd.Run` will
never pick the wrong definition.

`go-where` differs from [`godef`](https://github.com/rogpeppe/godef)
largely in that we do not assume input is a byte position in a Go
source file.


## Emacs and org-mode integration

The included [go-where.el](go-where.el) elisp file defines a
`go-where` function that can be called interactively in Emacs, with
`M-x go-where`.

To integrate that with `org-mode`, you can do this:

```elisp
(org-add-link-type "go" 'go-where)
```

And now, in your `.org` files, a link like
`[[go:github.com/tv42/where#Ident]]` can be opened with mouse or `C-c
C-o`.


## Future prospects

If `go doc` learns to output file positions, `go-where` may be
deprecated in favor of `go doc -pos IMPORTPATH IDENT`. We'll still need
a `#` separator to support the single input string use case.
