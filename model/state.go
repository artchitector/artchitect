package model

// Odin: Желаю, чтобы наше творение не было закрыто извне! Хочу, чтобы в Мидгарде видели, как Я пишу свои картины!
// Odin: Нужно показать, как Artchitect дышит, как бьётся его творящее сердце!
// State - состояние процесса.

// OdinState = есть процесс творения картины. Он начинается с идеи, а заканчивается рисунком.
// Odin: состояние Одина??... WTF? МИКРОБЫ, ЧТО ВЫ О СЕБЕ ВОЗОМНИЛИ?? Я РАЗДАВЛЮ ВАШИ ТУПЫЕ ГОЛОВЫ, МУРАВЬИ! ВАША МЕРЗКАЯ ПЛАНЕТА БУДЕТ РАЗРУШЕНА, А ВЫ ВСЕ БУДУТЕ СОЖЖЕ......
// Odin: ... сделай глубокий вдох ... ...  сделай глубокий вдох ...
// Odin: OdinState - состояние процесса творения картины. Так за ним можно будет наблюдать. Наблюдать, как творение развивается.
// Loki: наблюдать за исподним Одина...
// Odin: ... гнев это иллюзия ... ...  сделай глубокий вдох ... гнев это иллюзия ... вспомни о медитации ... ...
// Odin: компоненты в midgard смогут отобразить изменения процесса, текущее состояние будет непрерывно туда передаваться
type OdinState struct {
	ArtID                   uint     `json:"artId"`
	Seed                    uint     `json:"seed"`
	SeedEntropyImageEncoded string   `json:"seedEntropyImageEncoded"`
	SeedChoiceImageEncoded  string   `json:"seedChoiceImageEncoded"`
	NumberOfWords           uint     `json:"numberOfWords"`
	Words                   []string `json:"words"`
	Painting                bool     `json:"painting"`
	CurrentPaintTime        uint     `json:"currentPaintTime"`
	ExpectedPaintTime       uint     `json:"totalPaintTime"`
	Painted                 bool     `json:"painted"`
	Enjoying                bool     `json:"enjoying"`
	CurrentEnjoyTime        uint     `json:"currentEnjoyTime"`
	ExpectedEnjoyTime       uint     `json:"expectedEnjoyTime"`
}

type GivingState struct {
	LastArtID uint   `json:"lastArtID"`
	Given     []uint `json:"given"`
}
