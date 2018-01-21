package nqo

import (
	"github.com/petermattis/opttoy/v4/opt"
	"fmt"
	"strings"
)

const (
	schemaIndex = 100
)

func encode(e opt.Expr) string {
	out := encodeImpl(e)
	return strings.Join(out, ",")
}

func encodeImpl(e opt.Expr) []string {
	var out []string
	for i := 0; i < e.ChildCount(); i++ {
		out = append(out, encodeImpl(e.Child(i))...)
	}

	if e.Operator() == opt.ConstOp {
		out = append(out, fmt.Sprintf("%v", e.Private()))
	} else if e.Operator() == opt.VariableOp {
		out = append(out, fmt.Sprintf("%d", e.Private().(opt.ColumnIndex) + schemaIndex))
	} else if e.Operator() == opt.ScanOp {
		out = append(out, fmt.Sprintf("%d", e.Private().(opt.TableIndex) + schemaIndex))
	}

	out = append(out, fmt.Sprintf("%d", e.Operator()))
	return out
}
