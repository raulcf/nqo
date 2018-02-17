package nqo

import (
	"context"
	"fmt"
	"testing"

	"path/filepath"

	"github.com/cockroachdb/cockroach/pkg/sql/opt/optbuilder"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/xform"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/testutils/datadriven"
)

func TestGenerateJoins(t *testing.T) {
	actual := GenerateJoins(8)
	expected := 8462
	if len(actual) != expected {
		t.Fatalf("expected %d joins but found %d", expected, len(actual))
	}
	//if !reflect.DeepEqual(actual, expected) {
	//	t.Fatalf("expected joins: %s\nbut found: %s", expected, actual)
	//}
}

func TestGenerateQueries(t *testing.T) {
	actual := GenerateQueries(3)
	expected := 112
	if len(actual) != expected {
		t.Fatalf("expected %d queries but found %d", expected, len(actual))
	}
	//if !reflect.DeepEqual(actual, expected) {
	//	t.Fatalf("expected queries: %s\nbut found: %s", expected, actual)
	//}
}

func TestPrintAllQueries(t *testing.T) {
	//printAllQueries(t, 8)
	printAllEncoded(t)
}

func printAllQueries(t *testing.T, max int) {
	actual := GenerateQueries(max)
	for _, query := range actual {
		fmt.Println(query)
	}
}

func printAllEncoded(t *testing.T) {
	actual := GenerateQueries(8)
	ctx := context.Background()
	catalog := NewTestCatalog()
	o := xform.NewOptimizer(catalog, xform.OptimizeNone)

	path := "testdata/tpch"
	t.Run(filepath.Base(path), func(t *testing.T) {
		datadriven.RunTest(t, path, func(d *datadriven.TestData) string {
			stmt, err := parser.ParseOne(d.Input)
			if err != nil {
				d.Fatalf(t, "%v", err)
			}

			switch d.Cmd {
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

	for _, query := range actual {
		stmt, err := parser.ParseOne(query)
		if err != nil {
			t.Fatalf("%v", err)
		}

		b := optbuilder.New(ctx, o.Factory(), stmt)
		root, props, err := b.Build()
		if err != nil {
			t.Fatalf("error: %v\n", err)
		}
		exprView := o.Optimize(root, props)
		e := encoder{}

		fmt.Println(e.encode(exprView))
	}
}
