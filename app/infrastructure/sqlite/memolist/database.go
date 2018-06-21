package memolist

import (
	"github.com/jinzhu/gorm"
	// import sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	domain "github.com/naokirin/slan-go/app/domain/memolist"
)

// Memo data
type Memo struct {
	gorm.Model
	Kind string
	User string
	Text string
}

var _ domain.Repository = (*Memo)(nil)

// GetID return id
func (m *Memo) GetID() uint {
	return m.ID
}

// GetUser returns user id
func (m *Memo) GetUser() string {
	return m.User
}

// GetText returns memo text
func (m *Memo) GetText() string {
	return m.Text
}

// GetKind returns memo kind
func (m *Memo) GetKind() string {
	return m.Kind
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
func (m *Memo) All(kind string, user string) []domain.Memo {
	db := connectDb()
	defer db.Close()
	all := []*Memo{}
	db.Where("user = ? AND kind = ?", user, kind).Order("created_at").Find(&all)
	result := []domain.Memo{}
	for _, c := range all {
		result = append(result, c)
	}
	return result
}

// DeleteAll delete all memo for user
func (m *Memo) DeleteAll(kind string, user string) {
	db := connectDb()
	defer db.Close()
	db.Where("user = ? AND kind = ?", user, kind).Delete(&Memo{})
}

// Add add memo
func (m *Memo) Add(kind string, user string, text string) {
	db := connectDb()
	defer db.Close()
	db.Create(&Memo{Kind: kind, User: user, Text: text})
}

// Delete delete memo
func (m *Memo) Delete(v domain.Memo) bool {
	db := connectDb()
	defer db.Close()
	db.Where("id = ?", v.GetID()).Delete(&Memo{})
	return true
}
