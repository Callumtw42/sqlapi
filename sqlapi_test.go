package sqlapi

import (
	"testing"

)

func Test_sqlConvertTypes(t *testing.T) {

	// MysqlConnect(
	// 	"root",
	// 	"0089fxcy?",
	// 	"localhost",
	// 	"3306",
	// 	"test",
	// )

	type args struct {
		rows *sql.Rows
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{"empty", args{"./empty.sql"}, []map[string]interface{}{}},
		{"populated", args{"./populated.sql"}, []map[string]interface{}{
			{"column_1": 1, "column_2": "apple"},
			{"column_1": 2, "column_2": "banana"},
			{"column_1": 3, "column_2": "cranberry"}}
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sqlConvertTypes(tt.args.rows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sqlConvertTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}
