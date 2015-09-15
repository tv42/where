package where_test

import (
	"go/build"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/tv42/where"
)

var (
	testDataDir  string
	testDataOnce sync.Once
)

func locateTestData(t testing.TB) string {
	testDataOnce.Do(func() {
		pkg, err := build.Import("github.com/tv42/where", ".", build.FindOnly)
		if err != nil {
			t.Fatalf("cannot locate test package: %v", err)
		}
		testDataDir = filepath.Join(pkg.Dir, "testdata")
	})
	return testDataDir
}

func buildWithTags(tags ...string) *build.Context {
	c := build.Default
	c.BuildTags = tags
	return &c
}

type testCase struct {
	buildCtx *build.Context
	ident    string
	want     string
}

func TestSimple(t *testing.T) {

	for idx, test := range []testCase{
		{nil, "sample-1#one", "src/sample-1/one.go:3:5"},
		{nil, "sample-1#two", "src/sample-1/two.go:4:2"},
		{nil, "sample-1#three", "src/sample-1/three.go:3:6"},
		{nil, "sample-1#four", "src/sample-1/four.go:3:6"},
		{nil, "sample-1#five.a", "src/sample-1/five.go:4:2"},
		{nil, "sample-1#five.a.b", "src/sample-1/five.go:8:2"},
		{nil, "sample-1#six.method", "src/sample-1/six.go:5:14"},
		{nil, "sample-1#six.ptrMethod", "src/sample-1/six.go:7:15"},
		{nil, "sample-1#seven.method", "src/sample-1/seven.go:7:17"},
		{nil, "sample-1#seven.a.method", "src/sample-1/seven.go:12:21"},
		{nil, "cgo-1#one", "src/cgo-1/cgo.go:5:6"},
		{nil, "tags-1#one", "src/tags-1/one_std.go:5:5"},
		{buildWithTags("xyzzy"), "tags-1#one", "src/tags-1/one_xyzzy.go:5:5"},
	} {
		buildCtxOrig := test.buildCtx
		if buildCtxOrig == nil {
			buildCtxOrig = &build.Default
		}
		buildCtx := *buildCtxOrig
		buildCtx.GOPATH = locateTestData(t)
		pos, err := where.Ident(test.ident, &buildCtx)
		if err != nil {
			t.Errorf("error on #%d: %q: %v", idx, test.ident, err)
			continue
		}
		if pos == nil {
			t.Errorf("nil pos on success #%d: %q", idx, test.ident)
			continue
		}
		p := *pos
		p.Filename = strings.TrimPrefix(p.Filename, buildCtx.GOPATH+"/")
		if g, e := p.String(), test.want; g != e {
			t.Errorf("wrong #%d %q: %q != %q", idx, test.ident, g, e)
		}
	}
}
