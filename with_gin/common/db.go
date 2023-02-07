package common

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

type Db struct {
	*gorm.DB
}

var DB *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./../gorm.db")
	if err != nil {
		fmt.Println("Error when trying to initialize DB: ", err)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(15)
	DB = db
	return db
}

func TestDBInit() *gorm.DB {
	test_db, err := gorm.Open("sqlite3", "./../gorm_test.db")
	if err != nil {
		fmt.Println("db err: (TestDBInit) ", err)
	}
	test_db.DB().SetMaxIdleConns(3)
	test_db.LogMode(true)
	DB = test_db
	return DB
}

// Delete the database after running testing cases.
func TestDBFree(test_db *gorm.DB) error {
	test_db.Close()
	err := os.Remove("./../gorm_test.db")
	return err
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
