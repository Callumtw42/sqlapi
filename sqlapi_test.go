package sqlapi

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/assert"
)

func TestSel(t *testing.T) {

	MysqlConnect(
		"root",
		"0089fxcy?",
		"localhost",
		"3306",
		"test",
	)

	type args struct {
		sqlPath string
	}
	tests := []struct {
		name     string
		args     args
		wantData []map[string]interface{}
	}{
		{"empty", args{"./empty.sql"}, []map[string]interface{}(nil)},
		{"populated", args{"./populated.sql"}, []map[string]interface{}{
			{"column_1": int64(1), "column_2": "apple"},
			{"column_1": int64(2), "column_2": "banana"},
			{"column_1": int64(3), "column_2": "cranberry"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData := Sel(tt.args.sqlPath)
			assert.EqualValues(t, gotData, tt.wantData, "", "")
		})
	}
}
