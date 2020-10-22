package goseeder

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// Seeder root seeder offering access to db connection and util functions
type Seeder struct {
	DB      *sql.DB
	context clientSeeder
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
	var names = ""

	flag.BoolVar(&seed, "gseed", seed, "goseeder - if set will seed")
	flag.StringVar(&env, "gsenv", "", "goseeder - env for which seeds to execute")
	flag.StringVar(&names, "gsnames", "", "goseeder - comma separated seeder names to run specific ones")
	flag.Parse()

	if !seed {
		clientMain()
		return
	}

	var seeders = make([]string, 0)
	if len(names) > 0 {
		seeders = strings.Split(names, ",")
	}

	execute(conProvider(), env, seeders...)
	os.Exit(0)
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
	r := regexp.MustCompile(`.*\.(?P<name>[a-zA-Z]+$)`)
	match := r.FindStringSubmatch(getFunctionName(seeder))

	seeders = append(seeders, clientSeeder{
		env:  env,
		name: match[1],
		cb:   seeder,
	})
}

// Execute will executes the given seeder method
func execute(db *sql.DB, env string, seedMethodNames ...string) {
	// Execute all seeders if no method name is given
	if len(seedMethodNames) == 0 {
		if env == "" {
			log.Println("Running all seeders...")
		} else {
			log.Printf("Running seeders for env %s...\n", env)
		}
		for _, seeder := range seeders {
			if env == "" || env == seeder.env {
				seed(&Seeder{
					DB:      db,
					context: seeder,
				})
			}
		}
		return
	}

	for _, seeder := range seeders {
		if _, r := findString(seedMethodNames, seeder.name); (env == "" || env == seeder.env) && r {
			seed(&Seeder{
				DB:      db,
				context: seeder,
			})
		}
	}
}

func seed(seeder *Seeder) {
	clientSeeder := seeder.context
	start := time.Now()
	log.Printf("[%s] started seeding...\n", clientSeeder.name)

	defer func() {
		if r := recover(); r != nil {
			log.Print(printError(fmt.Sprintf("[%s] seed failed: %+v\n", clientSeeder.name, r)))
		}
	}()

	clientSeeder.cb(*seeder)
	elapsed := time.Since(start)
	log.Printf("[%s] seeded successfully, duration %s\n", clientSeeder.name, elapsed)
}
