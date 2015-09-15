package where

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
)

// The default importer in go/importer works based on .a files, that's
// not very useful. We need to implement our own, also to get control
// of tags.

type myImporter struct {
	fset     *token.FileSet
	buildCtx *build.Context

	cache map[string]*types.Package
}

func newImporter(buildCtx *build.Context, fset *token.FileSet) types.Importer {
	return &myImporter{
		fset:     fset,
		buildCtx: buildCtx,
		cache:    make(map[string]*types.Package),
	}
}

var _ types.Importer = myImporter{}

func astFiles(pkg *ast.Package) []*ast.File {
	r := make([]*ast.File, 0, len(pkg.Files))
	for _, f := range pkg.Files {
		r = append(r, f)
	}
	return r
}

func (imp myImporter) Import(importPath string) (*types.Package, error) {
	if pkg, ok := imp.cache[importPath]; ok {
		return pkg, nil
	}
	pkgMeta, err := imp.buildCtx.Import(importPath, ".", 0)
	if err != nil {
		return nil, fmt.Errorf("cannot locate package: %v", err)
	}

	filter := func(fi os.FileInfo) bool {
		base := filepath.Base(fi.Name())
		for _, name := range pkgMeta.GoFiles {
			if name == base {
				return true
			}
		}
		for _, name := range pkgMeta.CgoFiles {
			if name == base {
				return true
			}
		}
		return false
	}
	pkgs, err := parser.ParseDir(imp.fset, pkgMeta.Dir, filter, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot parse package: %v", err)
	}
	pkgAST, ok := pkgs[pkgMeta.Name]
	if !ok {
		return nil, fmt.Errorf("package not found: no %q in %s", pkgMeta.Name, pkgMeta.Dir)
	}

	conf := types.Config{
		IgnoreFuncBodies: true,
		FakeImportC:      true,
		// TODO maybe set Error here, to be more lenient?
		Importer:                 imp,
		DisableUnusedImportCheck: true,
	}
	pkg, err := conf.Check(importPath, imp.fset, astFiles(pkgAST), nil)
	if err != nil {
		return nil, fmt.Errorf("cannot type check package: %v", err)
	}
	imp.cache[importPath] = pkg
	return pkg, nil
}
