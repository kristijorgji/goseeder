package goseeder

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
)

var dataPath = "db/seeds/data"

//SetDataPath this will allow you to specify a custom path where your seed data is located
func SetDataPath(path string) {
	dataPath = path
}

//FromJson inserts into a database table with name same as the filename all the json entries
func (s Seeder) FromJson(filename string) {
	s.fromJson(filename, filename)
}

//FromJsonIntoTable reads json from the given file and inserts into the given table name
func (s Seeder) FromJsonIntoTable(filename string, tableName string) {
	s.fromJson(filename, tableName)
}

func (s Seeder) fromJson(filename string, tableName string) {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", dataPath, filename))
	if err != nil {
		log.Fatal(err)
	}

	m := []map[string]interface{}{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		panic(err)
	}

	for _, e := range m {
		stmQuery, args := prepareStatement(tableName, e)
		stmt, _ := s.DB.Prepare(stmQuery.String())
		_, err := stmt.Exec(args...)
		if err != nil {
			panic(err)
		}
	}
}
