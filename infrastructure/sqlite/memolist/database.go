package memolist

import (
	"github.com/jinzhu/gorm"
	// import sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	domain "github.com/naokirin/slan-go/domain/memolist"
)

// Memo data
type Memo struct {
	gorm.Model
	User string
	Text string
}

// GetUser returns user id
func (m Memo) GetUser() string {
	return m.User
}

// GetText returns memo text
func (m Memo) GetText() string {
	return m.Text
}

func connectDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "db/memolist.db")
	if err != nil {
		panic("failed to connect database: memolist.db")
	}
	db.AutoMigrate(&Memo{})
	return db
}

// All returns memolist for user
func (m Memo) All(user string) []domain.Memo {
	db := connectDb()
	defer db.Close()
	all := []Memo{}
	db.Where("user = ?", user).Order("created_at").Find(&all)
	result := []domain.Memo{}
	for _, c := range all {
		result = append(result, c)
	}
	return result
}

// DeleteAll delete all memo for user
func (m Memo) DeleteAll(user string) {
	db := connectDb()
	defer db.Close()
	db.Where("user = ?", user).Delete(&Memo{})
}

// Add add memo
func (m Memo) Add(user string, text string) {
	db := connectDb()
	defer db.Close()
	db.Create(&Memo{User: user, Text: text})
}

// Delete delete memo
func (m Memo) Delete(user string, v domain.Memo) bool {
	db := connectDb()
	defer db.Close()
	d := Memo{User: v.GetUser(), Text: v.GetText()}
	db.Delete(&d)
	return true
}
