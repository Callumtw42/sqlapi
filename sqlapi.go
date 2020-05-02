package sqlapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var db *sql.DB

//MysqlConnect connects to sql server
func MysqlConnect(user string, password string, host string, port string, database string) *sql.DB {
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+database)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("MySql Connected")
	}
	return db
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
		case "INT":
			types[i] = new(int32)
		case "VARCHAR":
			types[i] = new(string)
		default:
			types[i] = new(string)
		}
	}

	return types

}

//Sel parses results of SELECT statement in the passed sql file into JSON array and sends as http response
func Sel(res http.ResponseWriter, req *http.Request, sqlPath string) {
	//get query results
	sql, err := ioutil.ReadFile(sqlPath)
	handle(err)

	rows, err := db.Query(string(sql))
	handle(err)
	defer rows.Close()

	//array to store JSON objects
	var data []map[string]interface{}

	//column names
	cols, _ := rows.Columns()

	for rows.Next() {

		//array of typed addresses to write results to
		types := sqlConvertTypes(rows)

		// Scan the result into the column pointers...
		err := rows.Scan(types...)
		handle(err)

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := types[i]
			m[colName] = val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		data = append(data, m)
	}

	//format response into json and send
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(data)

}
