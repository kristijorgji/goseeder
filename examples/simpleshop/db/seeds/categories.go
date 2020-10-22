package seeds

import (
	"fmt"
	"github.com/kristijorgji/goseeder"
	"simpleshop/util"
)

func categoriesSeeder(s goseeder.Seeder) {
	s.FromJson("categories")
}

func testCategoriesSeeder(s goseeder.Seeder) {
	for i := 0; i < 100; i++ {
		stmt, _ := s.DB.Prepare(`INSERT INTO categories(id, name) VALUES (?,?)`)
		_, err := stmt.Exec(util.RandomInt(1, int64(^uint16(0))), []byte(fmt.Sprintf(`{"en": "%s"}`, util.RandomString(7))))
		if err != nil {
			panic(err)
		}
	}
}
