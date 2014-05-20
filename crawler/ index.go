package main

import (
  "fmt"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB gorm.DB

// initialize gorm
func init() {
	var err error
	DB, err = gorm.Open("postgres", "user=gosearch password=gosearch dbname=gosearch sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
	}
}

// 
func Setup() {
  fmt.Println("Setting up database tables for GoSearch")
  DB.Exec("DROP TABLE pages;DROP TABLE words;DROP TABLE locations;")
  DB.AutoMigrate(Page{})
  DB.AutoMigrate(Word{})
  DB.AutoMigrate(Location{})
}


type Page struct {
	Id              int64
	Url             string `sql:"size:255;not null;unique"`
	Title           string `sql:"size:255"`
	CreatedAt       time.Time
  UpdatedAt       time.Time
}

type Word struct {
	Id              int64
	Stem            string `sql:"size:255;not null"`
}

type Location struct {
	Id              int64
	Position        int64
  WordId          int64
  PageId          int64
}
