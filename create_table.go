package nqo

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/coltypes"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/types"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
)

// CreateTable creates a test table from a parsed DDL statement and
// adds it to the catalog.
func (c *TestCatalog) CreateTable(stmt *tree.CreateTable) *sqlbase.TableDescriptor {
	tn, err := stmt.Table.Normalize()
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}

	tbl := c.AddTable(tn.Table())
	hasPrimaryKey := false
	for _, def := range stmt.Defs {
		switch def := def.(type) {
		case *tree.ColumnTableDef:
			if def.PrimaryKey {
				hasPrimaryKey = true
			}
			err := c.addColumn(def, tn.Table())
			if err != nil {
				panic(err)
			}
		case *tree.UniqueConstraintTableDef:
			if def.PrimaryKey {
				hasPrimaryKey = true
			}
		}
		// TODO(rytaft): In the future we will likely want to check for unique
		// constraints, indexes, and foreign key constraints to determine
		// nullability, uniqueness, etc.
	}

	// If there is no primary key, add the hidden rowid column.
	if !hasPrimaryKey {
		c.AddColumn(tn.Table(), "rowid", types.Int, false /* nullable */, true /* hidden */)
	}

	return tbl
}

func (c *TestCatalog) addColumn(def *tree.ColumnTableDef, tableName string) error {
	nullable := !def.PrimaryKey && def.Nullable.Nullability != tree.NotNull
	typ := coltypes.CastTargetToDatumType(def.Type)
	return c.AddColumn(tableName, string(def.Name), typ, nullable, false /* hidden */)
}
