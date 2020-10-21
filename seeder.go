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

// Seeder type
type Seeder struct {
	DB *sql.DB
}

type clientSeeder struct {
	env  string
	name string
	cb   func(s Seeder)
}

var seeders []clientSeeder

func WithSeeder(conProvider func() *sql.DB, clientMain func()) {
	var seed bool = false
	var env string = ""
	var names string = ""

	flag.BoolVar(&seed, "gseed", seed, "goseeder - if set will seed")
	flag.StringVar(&env, "gsenv", "", "goseeder - env for which seeds to execute")
	flag.StringVar(&names, "gsnames", "", "goseeder - comma separated seeder names to run specific ones")
	flag.Parse()

	if !seed {
		clientMain()
		return
	}

	var seeders []string = make([]string, 0)
	if len(names) > 0 {
		seeders = strings.Split(names, ",")
	}

	execute(conProvider(), env, seeders...)
	os.Exit(0)
}

func Register(seeder func(s Seeder)) {
	RegisterForEnv("", seeder)
}

func RegisterForTest(seeder func(s Seeder)) {
	RegisterForEnv("test", seeder)
}

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
	s := Seeder{db}

	// Execute all seeders if no method name is given
	if len(seedMethodNames) == 0 {
		if env == "" {
			log.Println("Running all seeders...")
		} else {
			log.Printf("Running seeders for env %s...\n", env)
		}
		for _, seeder := range seeders {
			if env == "" || env == seeder.env {
				seed(&s, seeder)
			}
		}
		return
	}

	for _, seeder := range seeders {
		if _, r := findString(seedMethodNames, seeder.name); (env == "" || env == seeder.env) && r {
			seed(&s, seeder)
		}
	}
}

func seed(rootSeeder *Seeder, seeder clientSeeder) {
	start := time.Now()
	log.Printf("[%s] started seeding...\n", seeder.name)

	defer func() {
		if r := recover(); r != nil {
			log.Print(printError(fmt.Sprintf("[%s] seed failed: %+v\n", seeder.name, r)))
		}
	}()

	seeder.cb(*rootSeeder)
	elapsed := time.Since(start)
	log.Printf("[%s] seeded successfully, duration %s\n", seeder.name, elapsed)
}
