package nqo

import (
	"fmt"
	"strings"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/xform"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
)

const (
	maxCols = 1600
)

// An encoder encodes a query plan into a format that can be understood
// by a neural network.
type encoder struct {
	// columns stores the encoding of each column in the query plan, indexed
	// by the optimizer's unique column index.
	columns map[opt.ColumnIndex]string
}

// encode converts the given query plan into a comma-separated string
// of integers by walking the query tree in a post-order traversal, and
// printing the integers that represent the operator and arguments
// at each node.
func (e *encoder) encode(ev xform.ExprView) string {
	e.columns = make(map[opt.ColumnIndex]string)
	out := e.encodeImpl(ev)
	return strings.Join(out, ",")
}

func (e *encoder) encodeImpl(ev xform.ExprView) []string {
	out := make([]string, 0, ev.ChildCount()+1)
	for i := 0; i < ev.ChildCount(); i++ {
		out = append(out, e.encodeImpl(ev.Child(i))...)
	}

	if ev.Operator() == opt.ConstOp {
		out = append(out, fmt.Sprintf("%v", ev.Private()))
	} else if ev.Operator() == opt.ScanOp {
		// save the column information
		tblIndex := ev.Private().(opt.TableIndex)
		table := ev.Metadata().Table(tblIndex).(*sqlbase.TableDescriptor)
		for i := range table.Columns {
			colIndex := ev.Metadata().TableColumn(tblIndex, i)
			colID := tree.ColumnID(table.Columns[i].ID)
			e.columns[colIndex] = fmt.Sprintf("%d", uint32(table.ID * maxCols) + uint32(colID))
		}

		out = append(out, fmt.Sprintf("%d", table.ID * maxCols))
	} else if ev.Operator() == opt.VariableOp {
		colIndex := ev.Private().(opt.ColumnIndex)
		out = append(out, e.columns[colIndex])
	}

	out = append(out, fmt.Sprintf("%d", ev.Operator()))
	return out
}
