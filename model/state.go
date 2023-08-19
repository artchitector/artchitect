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

// FriggState - процесс объединения единств
// Odin: Frigg объединяет единства сложным и продолжительным процессом, который Я хочу показать жителям Мидгарда.
// Odin: Я лично отбираю лидеров в единства, использую свои навыки предвиденья (понимания энтропии) и чувство вкуса.
// Frigg: опишем всё происходящее следующей рекурсивной структурой.
// Frigg: при объединении верхнеуровневого единства сначала будут объединены его дочерние единства, а затем оно само
// Frigg: для этого добавлена рекурсивность, чтобы отобразить весь процесс сверху донизу рекурсивно
type FriggState struct {
	Unity           *Unity   `json:"unity"`
	Children        []*Unity `json:"children"`
	TotalApplicants uint     `json:"totalApplicants"`
	TotalLeads      uint     `json:"totalLeads"`
	Leads           []uint   `json:"leads"` // текущие лидеры, выбранные к сборке коллажа

	CollageStarted    bool `json:"collageStarted"` // момент, когда уже собраны все дети, лидеры и создаётся коллаж
	CollageFinished   bool `json:"collageFinished"`
	CurrentEnjoyTime  uint `json:"currentEnjoyTime"`
	ExpectedEnjoyTime uint `json:"expectedEnjoyTime"`

	Subprocess *FriggState `json:"subprocess"` // текущий дочерний элемент, который в обработке (подпроцесс соединения)
}

func NewFriggState(unity Unity) *FriggState {
	return &FriggState{
		Unity:    &unity,
		Children: make([]*Unity, 0),
		Leads:    make([]uint, 0),
	}
}

func (fs *FriggState) Active() *FriggState {
	if fs.Subprocess != nil {
		return fs.Subprocess.Active()
	}
	return fs
}

func (fs *FriggState) AddSubprocess(unity Unity) {
	fs.Active().Subprocess = NewFriggState(unity)
}

func (fs *FriggState) ClearSubprocess() {
	if fs.Subprocess == nil {
		return
	}
	if fs.Subprocess.Subprocess == nil {
		fs.Subprocess = nil
		return
	}
	fs.Subprocess.ClearSubprocess()
}
