package nqo

// This file is home to TestEncode, which is used for testing functionality
// for encoding query plans into a format understood by the neural networks.
//
// Each testfile contains testcases of the form
//   <command>
//   <SQL statement or expression>
//   ----
//   <expected results>
//
// The supported commands are:
//
//  - build
//
//    Builds a memo structure from a SQL query and outputs a representation
//    of the "expression view" of the memo structure.
//
//  - encode
//
//    Builds a memo structure from a SQL query and outputs an encoded version
//    of the query in a format that can be used as input to a neural network.
//
//  - exec-ddl
//
//    Parses a CREATE TABLE statement, creates a test table, and adds the
//    table to the catalog.
//

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/opt/optbuilder"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/xform"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/testutils/datadriven"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
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
			ctx := context.Background()
			catalog := NewTestCatalog()

			datadriven.RunTest(t, path, func(d *datadriven.TestData) string {
				stmt, err := parser.ParseOne(d.Input)
				if err != nil {
					d.Fatalf(t, "%v", err)
				}

				switch d.Cmd {
				case "build":
					o := xform.NewOptimizer(catalog, xform.OptimizeNone)
					b := optbuilder.New(ctx, o.Factory(), stmt)
					root, props, err := b.Build()
					if err != nil {
						return fmt.Sprintf("error: %v\n", err)
					}
					exprView := o.Optimize(root, props)
					return exprView.String()

				case "encode":
					o := xform.NewOptimizer(catalog, xform.OptimizeNone)
					b := optbuilder.New(ctx, o.Factory(), stmt)
					root, props, err := b.Build()
					if err != nil {
						return fmt.Sprintf("error: %v\n", err)
					}
					exprView := o.Optimize(root, props)
					e := encoder{}
					return e.encode(exprView)

				case "exec-ddl":
					if stmt.StatementType() != tree.DDL {
						d.Fatalf(t, "statement type is not DDL: %v", stmt.StatementType())
					}

					switch stmt := stmt.(type) {
					case *tree.CreateTable:
						tbl := catalog.CreateTable(stmt)
						return fmt.Sprintf("%s%s", tbl.String(), "\n")

					default:
						d.Fatalf(t, "expected CREATE TABLE statement but found: %v", stmt)
						return ""
					}

				default:
					d.Fatalf(t, "unsupported command: %s", d.Cmd)
					return ""
				}
			})
		})
	}
}
