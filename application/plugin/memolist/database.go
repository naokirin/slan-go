package memolist

import "github.com/jinzhu/gorm"

func connectDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "db/memolist.db")
	if err != nil {
		panic("failed to connect database: memolist.db")
	}
	db.AutoMigrate(&memo{})
	return db
}

func all(user string) []memo {
	db := connectDb()
	defer db.Close()
	all := []memo{}
	db.Where("user = ?", user).Order("created_at").Find(&all)
	return all
}

func deleteAll(user string) {
	db := connectDb()
	defer db.Close()
	db.Where("user = ?", user).Delete(&memo{})
}

func add(user string, text string) {
	db := connectDb()
	defer db.Close()
	db.Create(&memo{User: user, Text: text})
}

func delete(user string, m memo) bool {
	db := connectDb()
	defer db.Close()
	db.Delete(&m)
	return true
}
