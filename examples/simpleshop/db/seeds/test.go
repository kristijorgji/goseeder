package seeds

import "github.com/kristijorgji/goseeder"

func init() {
	goseeder.RegisterForTest(testCategoriesSeeder)
	goseeder.RegisterForTest(testExercisesSeeder)
}
