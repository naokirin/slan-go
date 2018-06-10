package stampgacha

import (
	"time"

	"github.com/jinzhu/gorm"
	// import sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	domain "github.com/naokirin/slan-go/app/domain/stampgacha"
)

// StampGacha data
type StampGacha struct {
	gorm.Model
	User        string
	LottingTime time.Time
}

var _ domain.Repository = (*StampGacha)(nil)

// GetUser returns user id
func (s *StampGacha) GetUser() string {
	return s.User
}

// GetLottingTime returns last lotting time
func (s *StampGacha) GetLottingTime() time.Time {
	return s.LottingTime
}

func connectDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "db/stamp_gacha.db")
	if err != nil {
		panic("failed to connect database: stamp_gacha.db")
	}
	db.AutoMigrate(&StampGacha{})
	return db
}

// LastLottingTime returns last lotting time
func (s *StampGacha) LastLottingTime(user string) (time.Time, bool) {
	db := connectDb()
	defer db.Close()
	gacha := []*StampGacha{}
	db.Where("user = ?", user).Find(&gacha)
	if len(gacha) == 0 {
		return time.Time{}, false
	}
	return gacha[0].GetLottingTime(), true
}

// SaveLottingTime save lotting time
func (s *StampGacha) SaveLottingTime(user string) {
	db := connectDb()
	defer db.Close()
	gacha := []*StampGacha{}
	db.Where("user = ?", user).Find(&gacha)
	now := time.Now()
	if len(gacha) == 0 {
		db.Create(&StampGacha{User: user, LottingTime: now})
		return
	}
	gacha[0].LottingTime = now
	db.Save(gacha[0])
}
