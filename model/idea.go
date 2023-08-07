package model

// MaxSeed - Seed для ИИ. Может быть нарисовано столько разных вариантов картины (с одинаковыми словами)
const MaxSeed = 4294967295 // numpy accepts from 0 to 4294967295

// Word - Слово
// Odin: вначале было слово... Из таких слов Фрейя поймёт, что нарисовать.
type Word struct {
	Word    string      `json:"word"`
	Entropy EntropyPack `json:"entropy"` // Odin: из этой энтропии Тор придумал слово
}

// Idea - pantheon.Odin всезнающ, и заранее знает как создать картину, которую надо нарисовать. Это знание он передаёт в своей идее.
type Idea struct {
	Seed                 uint        // Odin: Freyja может нарисовать MaxSeed вариантов одной и той же картины. Тут конкретный, который задумал Я.
	SeedEntropy          EntropyPack // сохраняем энтропию, которую видел Odin в момент воспоминания seed-числа
	NumberOfWordsEntropy EntropyPack // сохраняем энтропию, которую видел Odin в момент воспоминания количества слов
	Words                []Word      // Odin: Это слова, которые будут составлять основу идеи картины. Пример: "brain,smile,by hidari,Archangel,Lucifer,sauron,sharp,fractal,Tanks,moon and other planets and stars,by stanley"
}

func (l Idea) ExtractWords() []string {
	words := make([]string, 0, len(l.Words))
	for _, w := range l.Words {
		words = append(words, w.Word)
	}
	return words
}
