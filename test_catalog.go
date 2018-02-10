package nqo

import (
	"context"

	"fmt"

	"github.com/cockroachdb/cockroach/pkg/sql/optbase"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/types"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
)

// AddColumn adds the given column to the catalog. A table
// with the given table name must already exist.
func (c *TestCatalog) AddColumn(tableName, name string, typ types.T, nullable, hidden bool) error {
	table, ok := c.tables[tableName]
	if !ok {
		return fmt.Errorf("table %s does not exist in the catalog", tableName)
	}

	// Skip ID 0.
	c.nextColumnID[tableName]++

	colTyp, err := sqlbase.DatumTypeToColumnType(typ)
	if err != nil {
		return err
	}

	colDesc := sqlbase.ColumnDescriptor{
		Name:     name,
		ID:       sqlbase.ColumnID(c.nextColumnID[tableName]),
		Type:     colTyp,
		Nullable: nullable,
		Hidden:   hidden,
	}

	table.Columns = append(table.Columns, colDesc)
	return nil
}

// TestCatalog implements the optbase.Catalog interface for testing purposes.
type TestCatalog struct {
	tables       map[string]*sqlbase.TableDescriptor
	nextColumnID map[string]int
	nextTableID  int
}

var _ optbase.Catalog = &TestCatalog{}

// NewTestCatalog creates a new empty instance of the test catalog.
func NewTestCatalog() *TestCatalog {
	return &TestCatalog{
		tables:       make(map[string]*sqlbase.TableDescriptor),
		nextColumnID: make(map[string]int),
	}
}

// FindTable is part of the optbase.Catalog interface.
func (c *TestCatalog) FindTable(ctx context.Context, name *tree.TableName) (optbase.Table, error) {
	return c.tables[name.Table()], nil
}

// AddTable adds the given table to the catalog.
func (c *TestCatalog) AddTable(name string) *sqlbase.TableDescriptor {
	// Skip ID 0.
	c.nextTableID++
	tableDesc := &sqlbase.TableDescriptor{Name: name, ID: sqlbase.ID(c.nextTableID)}
	c.tables[name] = tableDesc
	c.nextColumnID[name] = 0
	return tableDesc
}
