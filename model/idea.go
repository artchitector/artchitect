package model

import (
	"fmt"
	"strings"
	"time"
)

// MaxSeed - Seed для ИИ. Может быть нарисовано столько разных вариантов картины (с одинаковыми словами)
const MaxSeed = 4294967295 // numpy accepts from 0 to 4294967295
const MaxKeywords = 28     // не более 28 слов

// Word - Слово
// Odin: вначале было слово... Из таких слов Фрейя поймёт, что нарисовать.
type Word struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	IdeaID    uint      `json:"-"`
	CreatedAt time.Time `json:"createdAt"`

	Word    string      `json:"word"`
	Entropy EntropyPack `json:"entropy" gorm:"embedded;embeddedPrefix:entropy_"` // Odin: из этой энтропии Тор придумал слово
}

// Idea - pantheon.Odin всезнающ, и заранее знает как создать картину, которую надо нарисовать.
// Odin: я предвижу картину и с помощью моих воронов Huginn и Muninn я объясню в этой идее, что надо рисовать
type Idea struct {
	// Odin: ID используется тот же идентификатор, что и у Art
	ArtID     uint      `json:"id" gorm:"primarykey"` // Odin: идею стоит сохранять после того, как картина нарисована успешно
	CreatedAt time.Time `json:"createdAt"`

	Seed                 uint        // Odin: Freyja может нарисовать MaxSeed вариантов одной и той же картины. Тут конкретный, который задумал Я.
	SeedEntropy          EntropyPack `gorm:"embedded;embeddedPrefix:seed_"`       // сохраняем энтропию, которую видел Odin в момент воспоминания seed-числа
	NumberOfWordsEntropy EntropyPack `gorm:"embedded;embeddedPrefix:nmbrofwrds_"` // сохраняем энтропию, которую видел Odin в момент воспоминания количества слов
	Words                []Word      // Odin: Это слова, которые будут составлять основу идеи картины. Пример: "brain,smile,by hidari,Archangel,Lucifer,sauron,sharp,fractal,Tanks,moon and other planets and stars,by stanley"
}

func (l Idea) ExtractWords() []string {
	words := make([]string, 0, len(l.Words))
	for _, w := range l.Words {
		words = append(words, w.Word)
	}
	return words
}

func (l Idea) String() interface{} {
	words := l.ExtractWords()
	return fmt.Sprintf("IDEA S:%d W:%s", l.Seed, strings.Join(words, ","))
}
