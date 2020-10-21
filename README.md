# goseeder

#### Motivation
While golang is a great language and getting better, there are still a lot of pieces missing for developing fast and in accurate way to avoid repetitions.
 
I was searching for a go seeder similar to the one that Laravel/Lumen provides and could not find one.

Knowing that this is such an important key element of any big project for testing and seeding projects with dummy data I decided to create one myself and share.

#### Features
For now the library supports only MySql as a database driver for its utility functions like `FromJson` but is db agnostic for your custom seeders you can use any database that is supported by `sql.DB`

`goseeder`
1. Allows specifying seeds for different environments such as test,common (all envs) and more (flexible you can define them) 
2. Provides out of the box functions like `FromJson` to seed the table from json data and more data formats and drivers coming soon

# Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [License](#license)

## Installation

```sh
go get github.com/kristijorgji/goseeder
```

## Usage

To use goseeder you just need to wrap your main function with `WithSeeder` and specify your seeds under a local package named `seeds` in one folder with structure `db/seeds` (it is important if you want to use json or other data to follow this structure)

// TODO full example usage will be uploaded later in this repo


## License

goseeder is released under the MIT Licence. See the bundled LICENSE file for details.






