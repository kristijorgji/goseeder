# goseeder

![GitHub Workflow Status](https://github.com/kristijorgji/goseeder/workflows/CI/badge.svg)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/kristijorgji/goseeder?tab=doc)
[![codecov](https://img.shields.io/codecov/c/github/kristijorgji/goseeder/badge.svg)](https://codecov.io/gh/kristijorgji/goseeder)
[![Go Report Card](https://goreportcard.com/badge/kristijorgji/goseeder)](https://goreportcard.com/report/kristijorgji/goseeder)
[![Sourcegraph](https://sourcegraph.com/github.com/kristijorgji/goseeder/-/badge.svg)](https://sourcegraph.com/github.com/kristijorgji/goseeder?badge)

#### Motivation
Golang is a great language and getting better as community and frameworks, but there are still a lot of pieces missing for developing fast, accurate way and avoiding repetitions.
 
I was searching for a go seeder similar to the one that Laravel/Lumen provides and could not find one.

Knowing that this is such an important key element of any big project for testing and seeding projects with dummy data I decided to create one myself and share.

#### Features
For now the library supports only MySql as a database driver for its utility functions like `FromJson` provided by `Seeder` struct, but it is db agnostic for your custom seeders you can use any database that is supported by `sql.DB`

`goseeder`
1. It is designed for different kind of usages, for both programmatically or building into your exe and run via cli args
1. Allows specifying seeds for different environments such as predefined test and custom specified envs by the user
2. Allows specifying list (or single) seed name for execution
3. Allows having common seeds that execute for every env, unless specified not to do so with the respective cli or option
4. Provides out of the box functions like `(s Seeder) FromJson` to seed the table from json data and more data formats and drivers coming soon

# Table of Contents

- [Installation](#installation)
- [Usage method 1: Turn your executable into seedable via cli args](#usage-method-1-turn-your-executable-into-seedable-via-cli-args)
    - [1. Change Your Main Function](#1-change-your-main-function)
    - [2. Registering Your Seeds](#2-registering-your-seeds)
    - [3. Run Seeds Only For Specific Env](#3-run-seeds-only-for-specific-env)
    - [4. Run Seeds By Name](#4-run-seeds-by-name)
- [Usage method 2: Programmatically](#usage-method-2-programmatically)    
- [Summary Of Cli Args](#summary-of-cli-args)
- [License](#license)

## Installation

```sh
go get github.com/kristijorgji/goseeder
```

## Usage method 1: Turn your executable into seedable via cli args

[Please check examples/simpleshop](examples/simpleshop) for a full working separate go project that uses the seeder

Below I will explain once more all the steps needed to have goseeder up and running for your project.

### 1. Change Your Main Function

In order to give your executable seeding abilities and support for its command line arguments, the first thing we have to do is to wrap our main function
with the provided `goseeder.WithSeeder`

```go
func WithSeeder(conProvider func() *sql.DB, clientMain func())
```

The function requires as argument 
1. one function that returns a db connection necessary to seed
2. your original main function, which will get executed if no seed is requested

One such main file can look like below:

```go
// main.go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/kristijorgji/goseeder"
	"log"
	"net/url"
	"os"
	_ "simpleshop/db/seeds"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	goseeder.WithSeeder(connectToDbOrDie, func() {
		myMain()
	})
}

func myMain() {
	fmt.Println("Here you will execute whatever you were doing before using github.com/kristijorgji/goseeder like start your webserver etc")
}

func connectToDbOrDie() *sql.DB {
	dbDriver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		url.QueryEscape(dbPassword),
		dbHost,
		dbPort,
		dbName,
	)
	con, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	return con
}

```

### 2. Registering Your Seeds

Great! After step 1 our executable is able to run in seed mode or default mode.

Now we want to know how to register our custom seeds.

If you look at the imports in main file from step one, we might notice that we import
`_ "simpleshop/db/seeds"` even though we do not use them directly. 

This is mandatory because our seeds will get registered during package initialisation as we will see later.

The recommended project folder structure to work properly with `goseeder` is to have the following path for the seeders `db/seeds` and the package name to be `seeds`

Inside the folder you can add your seeders, for example lets seed some data into the `categories` table from one json file located at `db/seeds/data/categories.json`.

To do that we create our `categories.go` file at `db/seeds` folder:

```go
// db/seeds/categories.go
package seeds

import (
	"github.com/kristijorgji/goseeder"
)

func categoriesSeeder(s goseeder.Seeder) {
	s.FromJson("categories")
}

```

To use this seed, the last step is to register it.

Seeds can be registered as:
1. `common` seeds that run for all environments
2. for a specific environment like `test`, `yourcustomenv` (more in step 3)

We are going to create below a seed that runs for all environments, so we will not specify any env while registering it.

To do that we create in the `db/seeds` folder the file `common.go` that will register seeds that get always executed regardless of the environment:
```go
// db/seeds/common.go
package seeds

import "github.com/kristijorgji/goseeder"

func init() {
	goseeder.Register(categoriesSeeder)
}
```

We used `goseeder.Register` to register our seed function to run for all environments.

That is all for the basic usage of goseeder!!! 

Our function in categories.go file is now registered and ready to be used.

Now you can run
```
go run main.go --gseed
```

and it will run all your seeds against the provided db connection.

The framework will look for `categories.json` file in the path `db/seeds/data`, and insert all the entries there in a table named `categories` (inferred from the file name)

If you have a seed registered for another environment, for example a test seed, the framework instead will look for the json file at `db/seeds/data/test`

So the rule is it will always lookup in this pattern `db/seeds/data/[environment]/[specifiedFileName].[type]`

You can also give a seed a custom name, if you do not want the function name to be used by default.
You can register a seed in a fully flexible way like:

```go
// db/seeds/common.go
package seeds

import "github.com/kristijorgji/goseeder"

func init() {
    Registration{
		Name: "another_name_for_cat_seeder",
		Env:  "stage",
	}.Complete(categoriesSeeder)
}
```

### 3. Run Seeds Only For Specific Env

Many times we want to have seeds only for `test` environment, test purpose and want to avoid having thousand of randomly generated rows inserted into production database by mistake!

Or we just want to have granular control, to have separate data to populate our app/web in different way for `staging` `prod` `yourcustomenv` and so on. 

goseeder is designed to take care of this by using one of the following methods:

- `goseeder.RegisterForTest(seeder func(s Seeder)` - registers the specified seed for the env named `test`
- `goseeder.RegisterForEnv(env string, seeder func(s Seeder))` - will register your seeder to be executed only for the custom specified env 

Let's add to our previous categories.go seeder one seed function specific only for test env!
The file now will look like:

```go
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

```

Finally, lets create our registrator file for all test seeds same way as we did with `common.go`, we will create `test.go` now with content as below:

```go
// db/seeds/test.go
package seeds

import "github.com/kristijorgji/goseeder"

func init() {
	goseeder.RegisterForTest(testCategoriesSeeder)
}

```

That is all! 

Now if you run your app without specifying test env, only the common env seeders will run and you cannot mess by mistake production or other environments!

To run the test seeder above you have to run:

```bash
go run main.go --gseed --gsenv=test
```

This will run only the tests registered for the env `test` and the `common` seeds. A seed is known as common if it is registered without envrionment via `Register` method and has empty string env.

If you do not want common seeds to get executed, just specify the flag `--gs-skip-common`

The above call to run only seeds for test `env`, and ignore the common ones then would be:

```bash
go run main.go --gseed --gsenv=test --gs-skip-common
```

### 4. Run Seeds By Name

When we register a seed like shown in step 2, the seed name is the same as the function name, so our seed is called `categoriesSeeder` because that is the name of the function we register below

```go
func init() {
	goseeder.Register(categoriesSeeder)
}
```

This is important because we are totally flexible and can do cool things like execute only the specific seed functions that we want!

Let's assume that we have 100 seed functions, and want to execute only one of them which is named categoriesSeeder (that we registered above) and ignore all the other seeds.

Easy as this, just run:

```bash
go run main.go --gseeder --gsnames=categoriesSeeder
```

If you want to execute multiple seeds by specifying their names, just use comma separated value like `--gsnames=categoriesSeeder,someOtherSeeder`

## Usage method 2: Programmatically

`goseeder` is designed to fit all needs for being the best seeding tool.

That means that you might want to seed data before your unit tests programmatically without using cli args. 

That is straightforward to do with goseeder.
Let us assume we want to test our api that connects to some database in the package `api`,

The file `api/main_test.go` might look like below:

```go
//api/main_test.go
package api

import (
	"database/sql"
	_ "db/seeds" // please import your seeds package so they register or register programatically here too if you want before the seeder Execute is called
	"log"
	"os"
	"testing"
    "github.com/kristijorgji/goseeder"
)


func TestMain(m *testing.M) {
	con := db.ConnectToTestDb()
	SeedTestData(con)
	r := m.Run()
	os.Exit(r)
}

func SeedTestData(con *sql.DB) {
	log.Println("Seeding test database")
	goseeder.SetDataPath("../db/seeds/data")
	err := goseeder.Execute(con, goseeder.ForEnv("test"), goseeder.ShouldSkipCommon(true))
	if err != nil {
		log.Fatal("Seeding test data failed\n")
		os.Exit(-2)
	}
}
```

How nice is that ?!

After the execution you have a database with data sourcing only from your seeds registered for test env (test seeds) !!!
The above is production code used by one company, but you might need to adjust to your needs.

Another common use case is to want to execute programmatically the seeder because you don't want to turn your executable into seedable (you don't want to use method 1)).

Then again you can just create another file `myseeder.go` and inside it do your custom logic or handling of args then just execute
`goseeder.Execute`

Your `myseeder.go` might look like
```go
package main

import (
	_ "db/seeds"
	"github.com/kristijorgji/goseeder"
)

func main() {
	err := goseeder.Execute(connectToDbOrDie())
    if err != nil {
        log.Fatal("Seeding test data failed\n")
        os.Exit(-2)
    }
}

// your func here to connect to db 
// connectToDbOrDie

```

Then you have your server or app executable separate for example in `main.go` file, and the seeder functionality separated in `myseeder.go`

You can easily run your seeder `go run myseeder.go`, or build and run etc based on your requirements.

You can pass all the necessary options to the  `goSeeder.Execute` method.
If you want to execute seeders for a particular env only (skip common seeds) for example you do it like:
```go
goseeder.Execute(con, goseeder.ForEnv("test"), goseeder.ShouldSkipCommon(true))
```

These are the possible options you can give after the mandatory db connection:
- `ForEnv(env string)` - you can specify here the env for which you want to execute
- `ForSpecificSeeds(seedNames []string)` - just specify array of seed names you want to execute
- `ShouldSkipCommon(value bool)` - this option has effect only if also gsenv if set, then will not run the common seeds (seeds that do not have any env specified
- `ShouldSkipCommon(value bool)` - this option has effect only if also gsenv if set, then will not run the common seeds (seeds that do not have any env specified



## Summary Of Cli Args

You can always  run
```bash
run go run main.go --help
````

to see all the available arguments and their descriptions. 

For the current version the result is:

```bash
INR00009:simpleshop kristi.jorgji$ go run main.go --help
Usage of /var/folders/rd/2bkszcpx6xgcddpn7f3bhczn1m9fb7/T/go-build358407825/b001/exe/main:
  -gs-skip-common
        goseeder - this arg has effect only if also gsenv if set, then will not run the common seeds (seeds that do not have any env specified)
  -gseed
        goseeder - if set will seed
  -gsenv string
        goseeder - env for which seeds to execute
  -gsnames string
        goseeder - comma separated seeder names to run specific ones
```

## License

goseeder is released under the MIT Licence. See the bundled LICENSE file for details.






