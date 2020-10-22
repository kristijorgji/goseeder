package seeds

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kristijorgji/goseeder"
	"simpleshop/util"
)

func productsSeeder(s goseeder.Seeder) {
	s.FromJson("products")
}

func testProductsSeeder(s goseeder.Seeder) {
	s.FromJson("products")
	for i := 0; i < 100; i++ {
		stmt, _ := s.DB.Prepare(`INSERT INTO products(id, name) VALUES (?, ?)`)
		_, err := stmt.Exec(
			uuid.New().String(),
			[]byte(fmt.Sprintf(`{"en": "%s"}`, util.RandomString(7))),
		)
		if err != nil {
			panic(err)
		}
	}
}
