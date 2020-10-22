package seeds

import "github.com/kristijorgji/goseeder"

func init() {
	goseeder.Register(categoriesSeeder)
	goseeder.Register(productsSeeder)
}
