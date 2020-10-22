package goseeder

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
)

//FromJson inserts into a database table with name same as the filename all the json entries
func FromJson(s Seeder, filename string) {
	content, err := ioutil.ReadFile(fmt.Sprintf("db/seeds/data/%s.json", filename))
	if err != nil {
		log.Fatal(err)
	}

	m := []map[string]string{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		panic(err)
	}

	for _, e := range m {
		stmQuery, args := prepareStatement(filename, e)
		stmt, _ := s.DB.Prepare(stmQuery.String())
		_, err := stmt.Exec(args...)
		if err != nil {
			panic(err)
		}
	}
}
