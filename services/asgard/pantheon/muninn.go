package pantheon

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
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
	return m.oneOf(ctx, model.MaxSeed)
}

// RememberNumberOfWords - Odin: сколько ключевых слов я назову, вспомнив эту картину
func (m *Muninn) RememberNumberOfWords(ctx context.Context) (uint, model.EntropyPack, error) {
	v, p, e := m.oneOf(ctx, model.MaxKeywords)
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
		index, pack, err := m.oneOf(ctx, count)
		if err != nil {
			return model.Word{}, errors.Wrap(err, "[muninn] ВЫБОР НЕ ОСУЩЕСТВЛЁН")
		}
		return model.Word{
			Word:    dict[index],
			Entropy: pack,
		}, nil
	}
}

// RememberArtNo - выбрать одну из картин в интервале номеров между min и max
func (m *Muninn) RememberArtNo(ctx context.Context, min uint, max uint) (uint, model.EntropyPack, error) {
	// Muninn: oneOf выбирает в интервале [0...99) (при limit=100). это безопасно для массивов
	// В выборе ID-картины логика иная, так как может выпасть и min, и max, и всё между ними
	// Допустим min=100, max=210. Всего может быть [100...210] элементов включительно
	// их всего 111 в этом множестве: от 100й до 199й - 100 элементов, от 200й до 210й - 11 элементов
	// значит в oneOf надо будет передать limit=111, и он вернёт число от 0 до 110. Это отступ.
	// Когда отступ прибавляется к min, то получится конкретный номер картины.
	// при выборе первой картины min+offset=100+0=100, а при выборе последней min+offset=100+110=210.
	// Так всё сойдётся.
	offset, ep, err := m.oneOf(ctx, max-min+1)
	if err != nil {
		return 0, model.EntropyPack{}, errors.Wrapf(err, "[muninn] ОШИБКА ВЫБОРА oneOf В ПРЕДЕЛАХ MIN=%d MAX=%d", min, max)
	}

	result := offset + min
	// TODO убрать это для отладки
	if result == min {
		log.Warn().Msgf("[muninn] ВЫПАЛ ПЕРВЫЙ ЭЛЕМЕНТ %d (min:%d, max:%d)", result, min, max)
	} else if result == max {
		log.Warn().Msgf("[muninn] ВЫПАЛ ПОСЛЕДНИЙ ЭЛЕМЕНТ %d (min:%d, max:%d)", result, min, max)
	}

	return result, ep, nil
}

// RememberUnityLead - Odin выбирает лидера, который будет составлять "лицо" единства (картина попадёт в коллаж)
func (m *Muninn) RememberUnityLead(ctx context.Context, ids []uint) (uint, model.EntropyPack, error) {
	if len(ids) == 0 {
		return 0,
			model.EntropyPack{},
			errors.Errorf("[muninn] НЕВОЗМОЖНО ВЫБРАТЬ ИЗ ПУСТОГО МАССИВА ЛИДЕРОВ")
	}
	if i, ep, err := m.oneOf(ctx, uint(len(ids))); err != nil {
		return 0,
			model.EntropyPack{},
			errors.Wrapf(err, "[muninn] ОШИБКА ВЫБОРА ЛИДЕРА ЕДИНСТВА ИЗ МАССИВА %d-ЭЛЕМЕНТОВ", len(ids))
	} else {
		return ids[i], ep, nil
	}
}

// oneOf - Muninn: если limit=100, то могут возвращаться числа от 0 до 99 (безопасно для выбора из массивов)
func (m *Muninn) oneOf(ctx context.Context, limit uint) (uint, model.EntropyPack, error) {
	pack, err := m.huginn.GetNextEntropy(ctx)
	if err != nil {
		return 0, model.EntropyPack{}, errors.Wrap(err, "[muninn] ХУГИН, ГДЕ ЭНТРОПИЯ?")
	}
	result := uint(math.Round(pack.Choice.FloatValue * float64(limit-1)))
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
