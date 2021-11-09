package inmem

import (
	"context"
	"fmt"
	"sync"
)

// An attempt to make a very generic inmem DB. To use, register desired tables and columns when NewDB is called
// EXAMPLE:
// func buildDB() *inmem.DB {
// 	tables := []inmem.Table{
// 		{
// 			Name:    "imports",
// 			Columns: []string{"csid", "id"},
// 		},
// 		{
// 			Name:    "profiles",
// 			Columns: []string{"id"},
// 		},
// 	}
//
//  return inmem.NewDB(tables)
// }

// DB ...
type DB struct {
	mu     sync.RWMutex
	tables map[string]*table
}

// NewDB ...
func NewDB(tables []Table) *DB {
	t := make(map[string]*table)
	for _, tbl := range tables {
		t[tbl.Name] = newTable(tbl.Columns...)
	}

	return &DB{
		tables: t,
	}
}

// Table is the exported representation of a table
type Table struct {
	Name    string
	Columns []string
}

// Get ...
func (db *DB) Get(ctx context.Context, table string, whereCol string, id string) ([][]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	tbl, found := db.tables[table]
	if !found {
		return nil, fmt.Errorf("table %q not found", table)
	}

	return tbl.get(val(id), colName(whereCol))
}

func (t *table) get(id val, whereCol colName) ([][]byte, error) {
	columnVals, found := t.rows[whereCol]
	if !found {
		return nil, fmt.Errorf("column %q not found", whereCol)
	}

	rowNums := columnVals[id]

	toReturn := make([][]byte, len(rowNums))
	for i, d := range rowNums {
		toReturn[i] = t.rowData[d]
	}

	return toReturn, nil
}

// Insert ...
func (db *DB) Insert(ctx context.Context, table string, cols []string, vals []string, data []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(cols) != len(vals) {
		return fmt.Errorf("length of cols must mach vals")
	}

	tbl, found := db.tables[table]
	if !found {
		return fmt.Errorf("table %q not found", tbl)
	}

	tbl.insert(row{
		cols: cols,
		vals: vals,
		data: data,
	})
	return nil
}

func (t *table) insert(r row) {
	t.rowData = append(t.rowData, r.data)
	rowNum := len(t.rowData) - 1

	for i, col := range r.cols {
		c, found := t.rows[colName(col)]
		if !found {
			c = make(map[val][]int)
			t.rows[colName(col)] = c
		}

		c[val(r.vals[i])] = append(c[val(r.vals[i])], rowNum)
	}

}

// Update ...
func (db *DB) Update(ctx context.Context, table string, col string, val string, data []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if col == "" || val == "" {
		return fmt.Errorf("column and value must be provided")
	}

	tbl, found := db.tables[table]
	if !found {
		return fmt.Errorf("table %q not found", tbl)
	}
	return tbl.update(col, val, data)
}

func (t *table) update(c, v string, d []byte) error {

	col, found := t.rows[colName(c)]
	if !found {
		return fmt.Errorf("column %q not found", c)
	}

	rowNums, found := col[val(v)]
	if !found {
		return fmt.Errorf("val %q not found", v)
	}

	for _, rowNum := range rowNums {
		t.rowData[rowNum] = d
	}
	return nil
}

type table struct {
	rows    map[colName]map[val][]int
	rowData [][]byte
}

func newTable(columns ...string) *table {
	r := make(map[colName]map[val][]int)
	for _, c := range columns {
		r[colName(c)] = make(map[val][]int)
	}
	return &table{
		rows: r,
	}
}

type row struct {
	cols []string
	vals []string
	data []byte
}

type val string
type colName string
