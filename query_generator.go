package nqo

import (
	"bytes"

	"github.com/cockroachdb/cockroach/pkg/util"
)

type table int

const (
	nation table = iota
	region
	part
	supplier
	partsupp
	customer
	orders
	lineitem
	numTables
)

var tableNames = [...]string{
	nation:   "nation",
	region:   "region",
	part:     "part",
	supplier: "supplier",
	partsupp: "partsupp",
	customer: "customer",
	orders:   "orders",
	lineitem: "lineitem",
}

type joinPath struct {
	pk   table
	fk   table
	cond string
}

type tabSet = util.FastIntSet

type joinCondKey struct {
	t1 table
	t2 table
}

var paths = []joinPath{
	{region, nation, "r_regionkey = n_regionkey"},
	{nation, supplier, "n_nationkey = s_nationkey"},
	{nation, customer, "n_nationkey = c_nationkey"},
	{part, partsupp, "p_partkey = ps_partkey"},
	{part, lineitem, "p_partkey = l_partkey"},
	{supplier, partsupp, "s_suppkey = ps_suppkey"},
	{supplier, lineitem, "s_suppkey = l_suppkey"},
	{partsupp, lineitem, "ps_partkey = l_partkey AND l_suppkey = ps_suppkey"},
	{customer, orders, "c_custkey = o_custkey"},
	{orders, lineitem, "o_orderkey = l_orderkey"},
}

var joinsWith map[table]*tabSet
var joinCondition map[joinCondKey]string

func makeJoinsWith() {
	joinsWith = make(map[table]*tabSet, numTables)
	for i := table(0); i < numTables; i++ {
		joinsWith[i] = &tabSet{}
	}

	joinCondition = make(map[joinCondKey]string, 2*len(paths))

	for i := range paths {
		path := &paths[i]

		joinsWith[path.pk].Add(int(path.fk))
		joinsWith[path.fk].Add(int(path.pk))

		joinCondition[joinCondKey{path.pk, path.fk}] = path.cond
		joinCondition[joinCondKey{path.fk, path.pk}] = path.cond
	}
}

func generate(n int) [][]table {
	tables := tabSet{}
	for i := 0; i < int(numTables); i++ {
		tables.Add(i)
	}
	return generateImpl(n, tables, tabSet{}, nil)
}

func generateImpl(n int, tables tabSet, join tabSet, joinSlice []table) [][]table {
	if n <= 0 {
		return [][]table{joinSlice}
	}

	var out [][]table
	tables.ForEach(func(i int) {
		if !join.Contains(i) {
			newJoin := join.Copy()
			newJoin.Add(i)
			newTables := joinsWith[table(i)].Difference(newJoin)
			newJoinSlice := make([]table, len(joinSlice))
			copy(newJoinSlice, joinSlice)
			newJoinSlice = append(newJoinSlice, table(i))
			out = append(out, generateImpl(n-1, newTables, newJoin, newJoinSlice)...)
		}
	})
	return out
}

func GenerateJoins(max int) [][]table {
	makeJoinsWith()
	var out [][]table
	for i := 1; i <= max; i++ {
		out = append(out, generate(i)...)
	}
	return out
}

func GenerateQueries(max int) []string {
	joins := GenerateJoins(max)
	var queries []string
	for _, join := range joins {
		var query bytes.Buffer
		query.WriteString("SELECT * FROM ")
		if len(join) > 0 {
			query.WriteString(tableNames[join[0]])
		}
		for i := 1; i < len(join); i++ {
			query.WriteString(" JOIN ")
			query.WriteString(tableNames[join[i]])
			query.WriteString(" ON ")
			query.WriteString(joinCondition[joinCondKey{join[i-1], join[i]}])
		}
		queries = append(queries, query.String())
	}
	return queries
}
