package portals

import (
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
