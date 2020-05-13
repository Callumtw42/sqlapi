package sqlapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

//DB database
var DB *sql.DB

//MysqlConnect connects to sql server
func MysqlConnect(user string, password string, host string, port string, database string) {
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+database)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("MySql Connected")
	}
	DB = db
}

func handle(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func sqlConvertTypes(rows *sql.Rows) []interface{} {

	rowInfo, _ := rows.ColumnTypes()
	types := make([]interface{}, len(rowInfo))

	for i := range rowInfo {
		switch rowInfo[i].DatabaseTypeName() {
		// case "INT", "BIGINT":
		// 	types[i] = new(sql.NullInt64)
		case "DOUBLE", "FLOAT", "INT", "BIGINT":
			types[i] = new(sql.NullFloat64)
		default:
			types[i] = new(sql.NullString)
		}
	}

	return types

}

//Sel parses results of SELECT statement in the passed sql file into JSON array and sends as http response
func Sel(sqlPath string) (data []map[string]interface{}) {
	//get query results
	query, err := ioutil.ReadFile(sqlPath)
	handle(err)

	rows, err := DB.Query(string(query))
	handle(err)
	defer rows.Close()

	//column names
	cols, _ := rows.Columns()

	//array of typed addresses to write results to
	addresses := sqlConvertTypes(rows)

	for rows.Next() {

		// Scan the result into the column pointers...
		err := rows.Scan(addresses...)
		handle(err)

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := reflect.ValueOf(addresses[i]).Elem().Field(0).Interface()
			m[colName] = val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		data = append(data, m)
	}

	return data

}

//Run executes SQL
func Run(sql string) {
	rows, err := DB.Query(string(sql))
	handle(err)
	defer rows.Close()
}

//JSONEncode formats response into json and send
func JSONEncode(res http.ResponseWriter, data []map[string]interface{}) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(data)
}
