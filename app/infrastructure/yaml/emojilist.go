package yaml

import (
	"sync"

	"github.com/naokirin/slan-go/app/domain/stampgacha"
)

// EmojiListRepository is emoji list
type EmojiListRepository struct {
	emojiList []string
}

var emojiListRepository *EmojiListRepository
var emojiListOnce sync.Once

// GetEmojiListRepository returns emojiList
func (r *EmojiListRepository) GetEmojiListRepository(path string) stampgacha.ConfigRepository {
	emojiListOnce.Do(func() {
		data, err := ParseFromFileToArray(path)
		if err != nil {
			return
		}
		list := make([]string, 0, len(data))
		for _, d := range data {
			list = append(list, d.(string))
		}
		emojiListRepository = &EmojiListRepository{
			emojiList: list,
		}
	})
	return emojiListRepository
}

// GetEmojiList returns emoji list
func (r *EmojiListRepository) GetEmojiList() []string {
	return r.emojiList
}
