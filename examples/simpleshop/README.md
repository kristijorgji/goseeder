# simpleshop example

This is one project example of using github.com/kristijorgji/goseeder

To have this example project up and running follow these steps:
1. you need to only create the database, just by executing the content of
`db/migrations/000001_init_schema.up.sql` against your database server
2. Then you need to create one .env file by copying the example .env.dist, and inside update the db credentials


Then in order to seed against run 
```bash 
go run main.go --gseed
```

and you will run all seeds in this example project.
 
If you do not want to seed, but run the original main logic of the project run without the gseed argument like:
 
 ```bash
 go run main.go
 ``` 


For full documentation of usage check the main documentation of github.com/kristijorgji/goseeder