package goseeder

//ConfigOption option to configure the seeder
type ConfigOption = func(config *config)

type config struct {
	env             string
	seedMethodNames []string
	skipCommon      bool
}

//ForEnv specify the env for which the seeders run for
func ForEnv(env string) ConfigOption {
	return func(config *config) {
		config.env = env
	}
}

//ForSpecificSeeds array of seed names you want to specify for execution
func ForSpecificSeeds(seedNames []string) ConfigOption {
	return func(config *config) {
		config.seedMethodNames = seedNames
	}
}

//ShouldSkipCommon this option has effect only if also gsenv if set, then will not run the common seeds (seeds that do not have any env specified
func ShouldSkipCommon(value bool) ConfigOption {
	return func(config *config) {
		config.skipCommon = value
	}
}
