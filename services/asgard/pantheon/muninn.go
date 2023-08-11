package pantheon

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"math"
	"os"
)

// Muninn - "помнящий". Один из воронов Odin-а.
// Muninn: я память Odin-а Всеотца, я могу вспомнить всё, что ему нужно для придумывания идеи картины.
// Mininn: я буду использовать интерпретацию энтропии от Huginn, и дам Odin-у нужные ему воспоминания
// Odin: я существую вне времени, и заранее знаю, что и когда будет нарисовано. Мне просто нужно вспомнить, что именно в этот раз.
type Muninn struct {
	// Muninn спрашивает Huginn о том, что Odin видит в своём LostEye
	// Muninn пробрасывает цепочку вызовов в Huginn, а оттуда в глаз Odin-а - LostEye
	huginn       *Huginn
	dictionaries map[string][]string
}

func NewMuninn(huginn *Huginn) *Muninn {
	return &Muninn{huginn: huginn, dictionaries: make(map[string][]string)}
}

// RememberSeed
// Odin: я вижу прошлое и будущее, и знаю, какая картина будет написана следующая.
// Odin: сейчас я напрягу свой глаз, а вороны Хугин и Мунин подскажут мне номер картины, который определит
// Odin: финальный вариант самого изображения. StableDiffusion знает, как с этим номером писать конкретную картинку
// Loki: так вот получается, ты изначально знаешь все картины Artchitect и сколько их будет. Сколько их?
// Odin: лишь одному Odin-у известно, сколько их, так оно и останется.
func (m *Muninn) RememberSeed(ctx context.Context) (uint, model.EntropyPack, error) {
	return m.OneOf(ctx, model.MaxSeed)
}

// RememberNumberOfWords - Odin: сколько ключевых слов я назову, вспомнив эту картину
func (m *Muninn) RememberNumberOfWords(ctx context.Context) (uint, model.EntropyPack, error) {
	v, p, e := m.OneOf(ctx, model.MaxKeywords)
	// количество слов от 1 до 28
	return v + 1, p, e
}

// RememberWord - вспомнить слово. model.Word уже содержит в себе model.EntropyPack внутри, потому он и не возвращается
func (m *Muninn) RememberWord(ctx context.Context) (model.Word, error) {
	dict, err := m.getDictionary(ctx, model.Version1)
	if err != nil {
		return model.Word{}, errors.Wrap(err, "[muninn] ОХ, МОЯ ПАМЯТЬ ОТКАЗАЛА! ODIN ЗАБЕРИ МЕНЯ В ВАЛЬХАЛЛУ!")
		// Odin: не ругай себя, малыш. Это всё ПОГРОМисты.
	} else {
		count := uint(len(dict))
		index, pack, err := m.OneOf(ctx, count)
		if err != nil {
			return model.Word{}, errors.Wrap(err, "[muninn] ВЫБОР НЕ ОСУЩЕСТВЛЁН")
		}
		return model.Word{
			Word:    dict[index],
			Entropy: pack,
		}, nil
	}

}

// OneOf - Muninn: если maxval=100, то могут возвращаться числа от 0 до 99 (безопасно для выбора из массивов)
func (m *Muninn) OneOf(ctx context.Context, maxval uint) (uint, model.EntropyPack, error) {
	pack, err := m.huginn.GetNextEntropy(ctx)
	if err != nil {
		return 0, model.EntropyPack{}, errors.Wrap(err, "[muninn] ХУГИН, ГДЕ ЭНТРОПИЯ?")
	}
	result := uint(math.Round(pack.Choice.FloatValue * float64(maxval-1)))
	return result, pack, nil
}

// getDictionary
// Muninn: извлекаю и восстанавливаю свою память о всех словах (которые знает Artchitect, конечно же)
// Muninn: Huginn поможет узнать, какое по порядку слово надо извлечь, а я найду его в памяти по номеру
// Muninn: используется сложная земная технология "YAML-файл"
func (m *Muninn) getDictionary(ctx context.Context, version string) ([]string, error) {
	if d, ok := m.dictionaries[version]; ok {
		return d, nil
	}
	filename := fmt.Sprintf("files/dictionary/%s.yaml", version)
	yamlData, err := os.ReadFile(filename)
	if err != nil {
		return []string{}, errors.Wrapf(err, "[muninn] ОШИБКА ЧТЕНИЯ ФАЙЛА %s", filename)
	}
	var words []string
	if err := yaml.Unmarshal(yamlData, &words); err != nil {
		// Muninn: извините, роняю Asgard.
		log.Fatal().Err(err).Msgf("[muninn] ОШИБКА ПАРСИНГА ФАЙЛА %s", filename)
		return []string{}, errors.Wrapf(err, "[muninn] ОШИБКА ПАРСИНГА ФАЙЛА %s", filename)
	}
	m.dictionaries[version] = words
	return words, nil
}
