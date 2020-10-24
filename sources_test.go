package goseeder

//
//import (
//	"github.com/DATA-DOG/go-sqlmock"
//	"testing"
//)
//
//func init() {
//	dataPath = "goseeder_testdata/db/seeds/data"
//}
//
//func TestSeeder_FromJson(t *testing.T) {
//	db, mock, _ := sqlmock.New()
//	s := Seeder{
//		DB:      db,
//		context: clientSeeder{},
//	}
//
//	mock.ExpectPrepare("")
//	mock.ExpectExec("insert into categories (name, id) values (?, ?)").WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectPrepare("")
//	mock.ExpectExec("insert into categories (name, id) values (?, ?)").WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
//	s.FromJson("categories")
//}
