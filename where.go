// Package where can locate definitions in Go source.
//
// The format of a string being looked up is
//
//     IMPORTPATH[#IDENTIFIER[.FIELD]...[.METHOD]]
//
// Import paths are separated from the rest by '#'. The rest of the
// string is a series of identifiers naming top-level declarations,
// struct fields, and methods.
//
// In practise, you can think it as "what godoc.org would take".
//
// Example:
//
//     github.com/tv42/where#Ident
package where

import (
	"fmt"
	"go/build"
	"go/token"
	"go/types"
	"strings"
)

// Ident looks up a definition by its import path and identifier names.
//
// buildCtx can be used to control GOPATH, GOOS, GOARCH, build tags,
// and such. If buildCtx is nil, go/build.Default is used.
func Ident(lookup string, buildCtx *build.Context) (*token.Position, error) {
	if buildCtx == nil {
		buildCtx = &build.Default
	}

	importPath, lookup := lookup, ""
	if idx := strings.IndexByte(importPath, '#'); idx >= 0 {
		importPath, lookup = importPath[:idx], importPath[idx+1:]
	}

	idents := strings.Split(lookup, ".")

	fset := token.NewFileSet()
	imp := newImporter(buildCtx, fset)

	pkg, err := imp.Import(importPath)
	if err != nil {
		return nil, fmt.Errorf("cannot import path: %q: %v", importPath, err)
	}

	obj := pkg.Scope().Lookup(idents[0])
	if obj == nil {
		return nil, fmt.Errorf("identifier not found: %q", idents[0])
	}
	for _, name := range idents[1:] {
		obj, _, _ = types.LookupFieldOrMethod(obj.Type(), true, pkg, name)
		if obj == nil {
			return nil, fmt.Errorf("field or method not found: %q", name)
		}
	}
	p := fset.Position(obj.Pos())
	return &p, nil
}
