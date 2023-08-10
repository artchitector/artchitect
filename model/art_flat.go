package model

import (
	"strings"
	"time"
)

// FlatArt - более плоская и компактная структура картины лишь с основными данными в плоском виде. Для отправки в Мидгард
// эта структура в БД не сохраняется, отправляется лишь в браузер
type FlatArt struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Version   string    `json:"version"`

	IdeaSeed          uint     `json:"ideaSeed"`
	IdeaNumberOfWords uint     `json:"ideaNumberOfWords"`
	IdeaWords         []string `json:"ideaWords"`

	// В списках используется лишь одна основная картинка энтропии - энтропия породившая Seed-номер
	SeedEntropyEncoded string `json:"imageEntropyEncoded"` //base64 encoded png картинка
	SeedChoiceEncoded  string `json:"imageChoiceEncoded"`  //base64 encoded png картинка
}

func MakeFlatArts(arts []Art) []FlatArt {
	fArts := make([]FlatArt, 0, len(arts))
	for _, a := range arts {
		fArts = append(fArts, MakeFlatArt(a))
	}
	return fArts
}

func MakeFlatArt(art Art) FlatArt {
	return FlatArt{
		ID:                 art.ID,
		CreatedAt:          art.CreatedAt,
		Version:            art.Version,
		IdeaSeed:           art.Idea.Seed,
		IdeaNumberOfWords:  uint(len(art.Idea.Words)),
		IdeaWords:          strings.Split(art.Idea.WordsStr, ","),
		SeedEntropyEncoded: art.Idea.SeedEntropy.Entropy.ImageEncoded,
		SeedChoiceEncoded:  art.Idea.SeedEntropy.Choice.ImageEncoded,
	}
}
