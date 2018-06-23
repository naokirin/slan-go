package vote

import (
	"github.com/jinzhu/gorm"
	// import sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	domain "github.com/naokirin/slan-go/app/domain/vote"
)

var _ domain.Repository = (*Vote)(nil)

// Vote is vote data model
type Vote struct {
	gorm.Model
	Title     string
	Owner     string
	OwnerName string
	Timestamp string
	Channel   string
	Hash      string
	AllowDup  bool
	Choises   []Choise
}

// Choise is one of vote choise model
type Choise struct {
	gorm.Model
	Number      int
	Text        string
	VoteID      uint `sql:"type:integer REFERENCES votes(id)"`
	UserChoises []UserChoise
}

// UserChoise is user choise model
type UserChoise struct {
	gorm.Model
	User     string
	ChoiseID uint `sql:"type:integer REFERENCES choises(id)"`
}

func (v *Vote) GetTitle() string {
	return v.Title
}

func (v *Vote) GetOwner() string {
	return v.Owner
}
func (v *Vote) GetTimestamp() string {
	return v.Timestamp
}
func (v *Vote) GetChannel() string {
	return v.Channel
}
func (v *Vote) GetOwnerName() string {
	return v.OwnerName
}
func (v *Vote) GetHash() string {
	return v.Hash
}
func (v *Vote) GetAllowDup() bool {
	return v.AllowDup
}

func (v *Vote) GetChoiseTexts() []string {
	result := make([]string, len(v.Choises))
	for _, c := range v.Choises {
		result[c.Number-1] = c.Text
	}
	return result
}

func (v *Vote) GetUserChoises() [][]string {
	result := make([][]string, len(v.Choises))
	for i := range result {
		result[i] = make([]string, 0)
	}
	for _, c := range v.Choises {
		for _, u := range c.UserChoises {
			result[c.Number-1] = append(result[c.Number-1], u.User)
		}
	}
	return result
}

func connectDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "db/vote.db")
	if err != nil {
		panic("failed to connect database: vote.db")
	}

	db.AutoMigrate(&Vote{})
	db.AutoMigrate(&Choise{})
	db.AutoMigrate(&UserChoise{})
	return db
}

// Create to db for vote
func (v *Vote) Create(hash, title, owner, ownerName, timestamp, channel string, allowDup bool, choises []string) {
	db := connectDb()
	defer db.Close()

	vote := &Vote{
		Title:     title,
		Owner:     owner,
		OwnerName: ownerName,
		Channel:   channel,
		Timestamp: timestamp,
		Hash:      hash,
		AllowDup:  allowDup,
	}
	db.Create(&vote)

	cs := make([]*Choise, len(choises))
	for i, c := range choises {
		cs[i] = &Choise{Number: i + 1, Text: c, VoteID: vote.ID}
		db.Create(cs[i])
	}
}

// AddUserChoise is registering user choise
func (v *Vote) AddUserChoise(user string, number int, timestamp string, channel string) {
	db := connectDb()
	defer db.Close()

	var vote Vote
	db.Where("timestamp = ? AND channel = ?", timestamp, channel).First(&vote)
	db.Where("vote_id = ?", vote.ID).Find(&vote.Choises)
	id := uint(0)
	for _, c := range vote.Choises {
		if c.Number == number {
			id = c.ID
			break
		}
	}
	userChoise := &UserChoise{
		User:     user,
		ChoiseID: id,
	}
	db.Create(&userChoise)
	for _, c := range vote.Choises {
		if c.Number == number {
			db.Model(&c).Association("UserChoises").Append(userChoise)
			break
		}
	}
}

// DeleteUserChoise is unregistering user choise
func (v *Vote) DeleteUserChoise(user string, number int, timestamp string, channel string) {
	db := connectDb()
	defer db.Close()

	var vote Vote
	db.Where("timestamp = ? AND channel = ?", timestamp, channel).First(&vote)
	db.Where("vote_id = ?", vote.ID).Find(&vote.Choises)
	for _, c := range vote.Choises {
		if c.Number == number {
			db.Where("user = ? AND choise_id = ?", user, c.ID).Delete(&UserChoise{})
			break
		}
	}
}

// Find returns registered user choises
func (v *Vote) Find(timestamp string, channel string) domain.Vote {
	db := connectDb()
	defer db.Close()

	var vote Vote
	db.First(&vote, "timestamp = ? AND channel = ?", timestamp, channel)
	db.Where("vote_id = ?", vote.ID).Find(&vote.Choises)
	for i := range vote.Choises {
		db.Where("choise_id = ?", vote.Choises[i].ID).Find(&vote.Choises[i].UserChoises)
	}
	return &vote
}

// FindByHash returns corresponding vote by hash
func (v *Vote) FindByHash(hash string) domain.Vote {
	db := connectDb()
	defer db.Close()

	var vote Vote
	db.First(&vote, "hash = ?", hash)
	db.Where("vote_id = ?", vote.ID).Find(&vote.Choises)
	for i := range vote.Choises {
		db.Where("choise_id = ?", vote.Choises[i].ID).Find(&vote.Choises[i].UserChoises)
	}
	return &vote
}
