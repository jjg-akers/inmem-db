package inmem_test

import (
	"context"
	"fmt"
	"os"

	"testing"

	"github.com/google/uuid"
	"github.com/jjg-akers/inmem-db/db/inmem"
	"github.com/jjg-akers/inmem-db/test"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(test.Coverage(m.Run(), 0.7, true))
}

func TestDB1_Get(t *testing.T) {

	csID := uuid.New().String()
	aggid := uuid.New().String()
	aggid2 := uuid.New().String()

	type args struct {
		ctx        context.Context
		tables     []inmem.Table
		tableToGet string
		getBy      string
		col        string
		aggID      string
	}

	testCases := []struct {
		name    string
		args    args
		setup   func(db *inmem.DB)
		want    [][]byte
		wantErr error
	}{
		{
			name: "should get empty import",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
				},
				tableToGet: "imports",
				getBy:      csID,
				col:        "csid",
				aggID:      uuid.New().String(),
			},
			setup: func(db *inmem.DB) {
				db.Insert(context.Background(), "imports", []string{"csid"}, []string{csID}, []byte{})
			},
			want: [][]byte{{}},
		},
		{
			name: "should get multiple imports by csid",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
				},
				tableToGet: "imports",
				getBy:      csID,
				col:        "csid",
				aggID:      uuid.New().String(),
			},
			setup: func(db *inmem.DB) {
				db.Insert(context.Background(), "imports", []string{"csid"}, []string{csID}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "processed", "fileName", "file1")))
				db.Insert(context.Background(), "imports", []string{"csid"}, []string{csID}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))
			},
			want: [][]byte{
				[]byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "processed", "fileName", "file1")),
				[]byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")),
			},
		},
		{
			name: "should get import by importID",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
				},
				tableToGet: "imports",
				getBy:      aggid,
				col:        "importID",
				aggID:      uuid.New().String(),
			},
			setup: func(db *inmem.DB) {
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "processed", "fileName", "file1")))
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid2}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))
			},
			want: [][]byte{
				[]byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "processed", "fileName", "file1")),
			},
		},
		{
			name: "should get import by importID",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
				},
				tableToGet: "imports",
				getBy:      aggid2,
				col:        "importID",
				aggID:      uuid.New().String(),
			},
			setup: func(db *inmem.DB) {
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "processed", "fileName", "file1")))
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid2}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))
			},
			want: [][]byte{
				[]byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")),
			},
		},
		{
			name: "should get correct data when mulitple tables are defined",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
					{
						Name: "users",
					},
				},
				tableToGet: "imports",
				getBy:      aggid2,
				col:        "importID",
				aggID:      uuid.New().String(),
			},
			setup: func(db *inmem.DB) {
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "processed", "fileName", "file1")))
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid2}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))
				db.Insert(context.Background(), "users", []string{"id"}, []string{aggid2}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))
			},
			want: [][]byte{
				[]byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")),
			},
		},
		{
			name: "should get correct data after an update - get unaffected row",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
					{
						Name: "users",
					},
				},
				tableToGet: "imports",
				getBy:      aggid2,
				col:        "importID",
			},
			setup: func(db *inmem.DB) {
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "processed", "fileName", "file1")))
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid2}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))

				db.Update(context.Background(), "imports", "importID", aggid, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "succeeded", "fileName", "file1")))

				db.Insert(context.Background(), "users", []string{"id"}, []string{aggid2}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))
			},
			want: [][]byte{
				[]byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")),
			},
		},
		{
			name: "should get correct data after an update - get updated row",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
					{
						Name: "users",
					},
				},
				tableToGet: "imports",
				getBy:      aggid,
				col:        "importID",
			},
			setup: func(db *inmem.DB) {
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "processed", "fileName", "file1")))
				db.Insert(context.Background(), "imports", []string{"csid", "importID"}, []string{csID, aggid2}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))

				db.Update(context.Background(), "imports", "importID", aggid, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "succeeded", "fileName", "file1")))

				db.Insert(context.Background(), "users", []string{"id"}, []string{aggid2}, []byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid2, "status", "succeeded", "fileName2", "file2")))
			},
			want: [][]byte{
				[]byte(fmt.Sprintf("{%q:%q, %q:%q, %q:%q}", "id", aggid, "status", "succeeded", "fileName", "file1")),
			},
		},
		{
			name: "should fail due to table not existing",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
					{
						Name: "users",
					},
				},
				tableToGet: "winky wonky",
				getBy:      aggid2,
				col:        "importID",
				aggID:      uuid.New().String(),
			},
			wantErr: fmt.Errorf("table %q not found", "winky wonky"),
		},
		{
			name: "should fail due to column not existing",
			args: args{
				ctx: context.Background(),
				tables: []inmem.Table{
					{
						Name: "imports",
					},
				},
				tableToGet: "imports",
				getBy:      aggid2,
				col:        "winky wonky",
				aggID:      uuid.New().String(),
			},
			wantErr: fmt.Errorf("column %q not found", "winky wonky"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db := inmem.NewDB(tc.args.tables)

			if tc.setup != nil {
				tc.setup(db)
			}

			gotImport, gotErr := db.Get(context.Background(), tc.args.tableToGet, tc.args.col, tc.args.getBy)
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, gotErr)
				return
			}

			if !assert.Nil(t, gotErr) {
				return
			}

			if !assert.Equal(t, len(tc.want), len(gotImport)) {
				return
			}

			for i, imp := range gotImport {
				assert.ElementsMatch(t, tc.want[i], imp)
			}
		})
	}
}
