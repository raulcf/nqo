package nqo

import (
	"flag"
	"path/filepath"
	"testing"

	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/opt/testutils"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/petermattis/opttoy/v4/build"
	"github.com/petermattis/opttoy/v4/cat"
	"github.com/petermattis/opttoy/v4/exec"
	"github.com/petermattis/opttoy/v4/opt"
)

var (
	testDataGlob = flag.String("d", "testdata/[^.]*", "test data glob")
)

func TestEncode(t *testing.T) {
	defer leaktest.AfterTest(t)()

	paths, err := filepath.Glob(*testDataGlob)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatalf("no testfiles found matching: %s", *testDataGlob)
	}

	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			catalog := cat.NewCatalog()

			testutils.RunDataDrivenTest(t, path, func(d *testutils.TestData) string {
				stmt, err := parser.ParseOne(d.Input)
				if err != nil {
					t.Fatal(err)
				}

				switch d.Cmd {
				case "exec":
					e := exec.NewEngine(catalog)
					return e.Execute(stmt)
				}

				p := opt.NewPlanner(catalog, 0 /* maxSteps */)
				b := build.NewBuilder(p.Factory(), stmt)
				root, required := b.Build()
				e := p.Optimize(root, required)

				if d.Cmd == "encode" {
					return fmt.Sprintf("%s\n", encode(e))
				}

				return e.String()
			})
		})
	}
}
