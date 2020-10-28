package goseeder

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// Seeder root seeder offering access to db connection and util functions
type Seeder struct {
	DB      *sql.DB
	context *clientSeeder
}

//Registration this allows for custom registration with full options available at once, like specifying custom seed name and env in one go. Then have to finish registration by calling Complete
type Registration struct {
	Name      string
	Env       string
	completed bool
}

//Complete this finished the registration of this registration instance. If you call a second time for same instance, error will be throw
func (r Registration) Complete(s func(seeder Seeder)) error {
	if r.completed {
		return errors.New("registration is already completed. You can use one registration only one time")
	}

	seeders = append(seeders, clientSeeder{
		env:  r.Env,
		name: r.Name,
		cb:   s,
	})
	r.completed = true

	return nil
}

type clientSeeder struct {
	env  string
	name string
	cb   func(s Seeder)
}

var seeders []clientSeeder

// WithSeeder It gives your main function seeding functions and provides the cli arguments
func WithSeeder(conProvider func() *sql.DB, clientMain func()) {
	var seed = false
	var env = ""
	var skipCommon = false
	var names = ""

	flag.BoolVar(&seed, "gseed", seed, "goseeder - if set will seed")
	flag.StringVar(&env, "gsenv", "", "goseeder - env for which seeds to execute")
	flag.StringVar(&names, "gsnames", "", "goseeder - comma separated seeder names to run specific ones")
	flag.BoolVar(&skipCommon, "gs-skip-common", skipCommon, "goseeder - this arg has effect only if also gsenv if set, then will not run the common seeds (seeds that do not have any env specified)")
	flag.Parse()

	if !seed {
		clientMain()
		return
	}

	var specifiedSeeders = make([]string, 0)
	if len(names) > 0 {
		specifiedSeeders = strings.Split(names, ",")
	}

	err := Execute(
		conProvider(),
		ForEnv(env),
		ForSpecificSeeds(specifiedSeeders),
		ShouldSkipCommon(skipCommon),
	)

	if err != nil {
		log.Panic(err)
	}
}

// Register the given seed function  as common to run for all environments
func Register(seeder func(s Seeder)) {
	RegisterForEnv("", seeder)
}

// RegisterForTest the given seed function for test environment
func RegisterForTest(seeder func(s Seeder)) {
	RegisterForEnv("test", seeder)
}

// RegisterForEnv the given seed function for a specific environment
func RegisterForEnv(env string, seeder func(s Seeder)) {
	r := regexp.MustCompile(`.*\.(?P<name>[a-zA-Z0-9]+$)`)
	match := r.FindStringSubmatch(getFunctionName(seeder))

	seeders = append(seeders, clientSeeder{
		env:  env,
		name: match[1],
		cb:   seeder,
	})
}

// Execute use this method for using this lib programmatically and executing
// seeder directly with full flexibility. Be sure first to have registered your seeders
func Execute(db *sql.DB, options ...ConfigOption) error {
	c := config{
		env:             "",
		seedMethodNames: []string{},
		skipCommon:      false,
	}
	for _, option := range options {
		option(&c)
	}

	// Execute all seeders if no method name is given
	if len(c.seedMethodNames) == 0 {
		if c.env == "" {
			printInfo(fmt.Sprintf("Running all seeders...\n\n"))
		} else {
			printInfo(fmt.Sprintf("Running all seeders for env %s and common seeds (without env)...\n\n", c.env))
		}
		for _, seeder := range seeders {
			if c.env == "" || c.env == seeder.env || (seeder.env == "" && c.skipCommon == false) {
				err := seed(&Seeder{
					DB:      db,
					context: &seeder,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	for _, seeder := range seeders {
		if _, r := findString(c.seedMethodNames, seeder.name); (c.env == "" || c.env == seeder.env) && r {
			err := seed(&Seeder{
				DB:      db,
				context: &seeder,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func seed(seeder *Seeder) (rerr error) {
	clientSeeder := seeder.context
	start := time.Now()
	printInfo(fmt.Sprintf("[%s] started seeding...\n", clientSeeder.name))

	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("[%s] seed failed: %+v\n", clientSeeder.name, r)
			printError(msg)
			rerr = errors.New(msg)
		}
	}()

	clientSeeder.cb(*seeder)
	elapsed := time.Since(start)
	printInfo(fmt.Sprintf("[%s] seeded successfully, duration %s\n\n", clientSeeder.name, elapsed))
	return nil
}
