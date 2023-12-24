package pantheon

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/services/asgard/pantheon/frigg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	unityRanks      = []uint{model.Unity100, model.Unity1K, model.Unity10K, model.Unity100K}
	updateIntervals = []uint{model.UpdateInterval100, model.UpdateInterval1K, model.UpdateInterval10K, model.UpdateInterval100K}
)

/*
Frigg - Фригг (др.-сканд. Frigg); — в германо-скандинавской мифологии жена Одина, верховная богиня.
Родоначальница рода асов. Покровительствует любви, браку, домашнему очагу, деторождению.
Является провидицей, которой известна судьба любого человека, но которая не делится этими знаниями ни с кем.
---
[ Odin рассказывает, зачем придумал единства]
Odin: Artchitect рисует больше тысячи картин каждый день, и сейчас нет простого способа их посмотреть.
Odin: в какой-то момент их число перевалит за сотни тысяч.
Odin: чтобы иметь некую навигацию по всему труду Artchitect - я придумал объединять картины в единства.
Odin: верховное единство - стотысячное, объединяет сто тысяч картин

	Пример: U0XXXXX (объедняет все картины с ведушим нулём, ID от 0 до 999999)

Odin: внутри него размещаются десятитысячные единства, тысячные единства, и затем сотни.

	Пример иерархии единств для картины #213650:
	- U2XXXXX (от 200000 до 299999) - model.Unity100K - стотысячное единство
	 - U21XXXX (от 210000 до 219999) - model.Unity10K - десятитысячное единство
	  - U213XXX (от 213000 до 213999) - model.Unity1K - тысячное единство
	   - U2136XX (от 213600 до 213699) - model.Unity100 - сотня картин

Odin: Я не буду заниматься всем сам лично, так как моя роль и так уже НАД ВСЕМ СУЩИМ.
Odin: И еще код Artchitect пытается иногда следовать принципам SOLID, поэтому...
Odin: Я попрошу свою супругу Фригг заняться объединением разрозненных частей в целое.
...
Frigg: Я как верховная богиня и покровительница домашнего очага,
Frigg: соберу всё отдельное в единое,
Frigg: как мать семейства собирает всех под одной крышей перед ужином,
Frigg: как Великая Мать собирает всё сущее в единую ойкумену непрерывно.
*/
type Frigg struct {
	collage                 *frigg.Collage // вспомогательная структура для разделения кода
	muninn                  *Muninn        // Мунин подсказывает, какую картину выбрать для коллажа
	heimdallr               *Heimdallr     // Хеймдалль поможет мне отправлять состояние единения в Мидгард
	unityPile               unityPile
	artPile                 artPile
	unificationEnjoyTimeSec uint
}

func NewFrigg(
	collage *frigg.Collage,
	muninn *Muninn,
	heimdallr *Heimdallr,
	unityPile unityPile,
	artPile artPile,
	unificationEnjoyTimeSec uint,
) *Frigg {
	return &Frigg{
		collage:                 collage,
		muninn:                  muninn,
		heimdallr:               heimdallr,
		unityPile:               unityPile,
		artPile:                 artPile,
		unificationEnjoyTimeSec: unificationEnjoyTimeSec,
	}
}

/*
ReunifyArtUnities - когда рисуется очередная картина, то это может вызвать каскадные обновления единств этой картины
Пример: сотенное единство обновляется каждые 10 картин, тысячное единство обновляется каждые 50 картин... (см константы в model.Unity)
Если на текущей картине единство завершено, то единство окончательно обновится до готового состояния и далее не будет пересобираться
напр. нарисована картина 099999, она последняя в единстве U0XXXXX,
и на ней полностью пересоберутся и закроются единства по порядку - U0999XX -> U099XXX -> U09XXXX -> U0XXXXX
*/
func (f *Frigg) ReunifyArtUnities(ctx context.Context, art model.Art) error {
	if err := f.createNonexistentUnities(ctx, art); err != nil {
		return errors.Wrapf(err, "[frigg] НЕВОЗМОЖНО ИНИЦИАЛИЗИРОВАТЬ ЕДИНСТВА ДЛЯ КАРТИНЫ #%d", art.ID)
	}
	unitiesToUpdate := make(map[string]struct{})
	nextArtID := art.ID + 1

	for _, rank := range unityRanks {
		if nextArtID%rank == 0 {
			// На следующей картине начнётся новое единство с рангом rank
			// Текущую сотню можно обновлять и закрывать
			unitiesToUpdate[art.GetUnityMask(rank)] = struct{}{}
		}
	}

	for _, interval := range updateIntervals {
		if nextArtID%interval == 0 {
			var rank uint
			switch interval {
			case model.UpdateInterval100:
				rank = model.Unity100
			case model.UpdateInterval1K:
				rank = model.Unity1K
			case model.UpdateInterval10K:
				rank = model.Unity10K
			case model.UpdateInterval100K:
				rank = model.Unity100K
			default:
				log.Fatal().Msgf("[frigg] КРИТИЧЕСКАЯ НЕИСПРАВНОСТЬ. НЕИЗВЕСТЕН ИНТЕРВАЛ %d. ВЫЗЫВАЮ ЗАВЕРШЕНИЕ АСГАРДА", interval)
			}

			unitiesToUpdate[art.GetUnityMask(rank)] = struct{}{}
		}
	}

	// Frigg: я собрала единства, которые нужно перестроить. Теперь я им проставлю статус model.UnityStateReunification
	// Frigg: и в следующем цикле Главного Цикла Творения они будут обновлены
	for mask := range unitiesToUpdate {
		unity, err := f.unityPile.Get(ctx, mask)
		if err != nil {
			return errors.Wrapf(err, "[frigg] ОШИБКА ПОИСКА ЕДИНСТВА %s", mask)
		}
		unity.State = model.UnityStateReunification
		unity, err = f.unityPile.Save(ctx, unity)
		if err != nil {
			return errors.Wrapf(err, "[frigg] ОШИБКА СБРОСА СТАТУСА ЕДИНСТВА %s НА %s", mask, model.UnityStateReunification)
		}
		log.Info().Msgf("[frigg] ЕДИНСТВО %s ПОДГОТОВЛЕНО К ОБНОВЛЕНИЮ", mask)
	}

	return nil
}

func (f *Frigg) HandleUnification(ctx context.Context) (worked bool, err error) {
	workUnity, err := f.unityPile.GetNextUnityForReunification(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // Frigg: для меня нет работы в этот раз.
	} else if err != nil {
		return false, errors.Wrap(err, "[frigg] ПРОБЛЕМЫ С ПОЛУЧЕНИЕМ ЕДИНСТВА ДЛЯ ОБЪЕДИНЕНИЯ")
	}

	workUnity, err = f.reunifyUnity(ctx, workUnity, nil)
	if err != nil {
		return false, errors.Wrapf(err, "[frigg] АВАРИЯ. ОБЪЕДИНЕНИЕ МНОЖЕСТВА %s", workUnity.Mask)
	}

	log.Info().Msgf(
		"[frigg] ОБЪЕДИНЕНИЕ МНОЖЕСТВА %s ЗАВЕРШЕНО. STATE=%s, VERSION=%d",
		workUnity.Mask,
		workUnity.State,
		workUnity.Version,
	)
	return true, nil
}

func (f *Frigg) createNonexistentUnities(ctx context.Context, art model.Art) error {
	for _, rank := range unityRanks {
		mask := art.GetUnityMask(rank)
		_, err := f.unityPile.Get(ctx, mask)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(err, "[frigg] СБОЙ ПОИСКА ЕДИНСТВА %s", mask)
		} else if err == nil {
			continue
		}
		min, max, err := f.getMinMaxID(mask)
		if err != nil {
			return errors.Wrapf(err, "[frigg] МАТЕМАТИЧЕСКИЙ СБОЙ MIN-MAX %s", mask)
		}

		state := model.UnityStateEmpty
		if art.ID >= max {
			state = model.UnityStateReunification
		}
		unity, err := f.unityPile.Create(ctx, mask, state, rank, min, max)
		if err != nil {
			return errors.Wrapf(err, "[frigg] СБОЙ СОЗДАНИЯ ЕДИНСТВА %s", mask)
		}
		log.Debug().Msgf("[frigg] СОЗДАНО ЕДИНСТВО %s", unity.Mask)

		if err := f.collage.MakeNSaveBlankCollage(ctx, unity); err != nil {
			return errors.Wrapf(err, "[frigg] СБОЙ СОХРАНЕНИЯ ПЕРВОЙ КАРТИНКИ ЕДИНСТВА %s", mask)
		}
	}

	return nil
}

func (f *Frigg) getMinMaxID(mask string) (min uint, max uint, err error) {
	var min64, max64 uint64
	min64, err = strconv.ParseUint(
		strings.ReplaceAll(strings.ReplaceAll(mask, "X", "0"), "U", ""),
		10,
		64,
	)
	if err != nil {
		return 0, 0, err
	}
	max64, err = strconv.ParseUint(
		strings.ReplaceAll(strings.ReplaceAll(mask, "X", "9"), "U", ""),
		10,
		64,
	)
	return uint(min64), uint(max64), nil
}

/*
reunifyUnity

Frigg: теперь я расскажу, как множества будут объединяться и что это за процесс.
Frigg: чтобы понять, что такое единство, нужно посетить страницу https://artchitect/unity и посмотреть единства.
Frigg: Множество с внешней точки зрения выглядит как набор карточек (коллаж из картин), входящий в это множество.
Frigg: Можно будет взглянуть на картинку множества U001XXX и увидеть сетку 5x5 из картин с номерами от 1000 до 1999
Frigg: Если углубиться в это множество, то там будут дочерние множества - сотни.
Frigg: так человек из Мидгарда сможет просматривать единства, и перемещаться по ним и вглубь, и вдоль.
Frigg: единства - удобный браузер картин Artchitect. Этих картин будут сотни тысяч, и нужна удобная форма для просмотра всего.

Frigg: картина единства начинает собираться почти сначала жизни единства, и в незаполненных единствах вместо картин будут
чёрные поля в местах пропуска.
Frigg: Каждую картину для коллажа единства выбирает лично Odin. Odin может указывать и на отсутствующие в данный момент картины,
и вместо них и будет видна чёрная пустая область. Чем больше заполнено единство картинами, тем меньше чёрных пропусков.
*/
func (f *Frigg) reunifyUnity(ctx context.Context, unity model.Unity, state *model.FriggState) (model.Unity, error) {
	log.Info().Msgf("[frigg] НАЧИНАЮ ОБЪЕДИНЯТЬ ЕДИНСТВО %s", unity.Mask)

	if state == nil {
		state = model.NewFriggState(unity)
		f.sendState(ctx, state)
	} else {
		state.AddSubprocess(unity)
	}

	// Frigg: при сборке единства сначала надо пройти по его дочерним единствам и объединить их
	if unity.Rank != model.Unity100 {
		// Frigg: внутри сотенного единства нет дочерних для объединения
		if children, err := f.unityPile.GetChildren(ctx, unity); err != nil {
			return model.Unity{}, errors.Wrapf(err, "[frigg] ОШИБКА ДОСТУПА К ДЕТЯМ ЕДИНСТВА %s", unity.Mask)
		} else {
			// Frigg: Сначала я набью шкурку содержимым
			for _, child := range children {
				ch := child
				state.Active().Children = append(state.Active().Children, &ch)
			}
			f.sendState(ctx, state)

			// Frigg: Сначала я набью шкурку содержимым, а только затем и её закончу
			for idx, child := range children {
				if child.State == model.UnityStateReunification {
					// рекурсивная сборка
					log.Info().Msgf("[frigg] НАЧИНАЮ ОБЪЕДИНЯТЬ ДОЧЕРНЕЕ ЕДИНСТВО %s", child.Mask)
					child, err = f.reunifyUnity(ctx, child, state)
					if err != nil {
						return model.Unity{}, errors.Wrapf(err, "[frigg] ОШИБКА ОБЪЕДИНЕНИЯ ДОЧЕРНЕГО ЕДИНСТВА %s", child.Mask)
					}

					ch := child
					state.Active().Children[idx] = &ch
				}
			}
		}
	}

	f.sendState(ctx, state)

	// Frigg: картинки, которые составляю коллаж, хранятся в массиве Leads у model.Unity.
	// Frigg: для сотенного единства картинки выбираются из всей сотни, но вот единства уровнем выше выбирают не любые картины,
	// Frigg: а лишь уже выбранные в дочерних сотнях. Это нужно, чтобы если человек увидел интересную картину у множества U01XXXX,
	// Frigg: то человек может зайти внутрь и найти дочернее единство с этой картинкой, и так дойти до сотни с той самой первой картинкой
	applicants, err := f.collectApplicants(ctx, unity)
	if err != nil {
		return model.Unity{}, errors.Wrapf(err, "[frigg] ОШИБКА СБОРА ПРЕТЕНДЕНТОВ ДЛЯ ЕДИНСТВА %s", unity.Mask)
	}

	state.Active().TotalApplicants = uint(len(applicants))
	f.sendState(ctx, state)

	leadsCount := 0
	switch unity.Rank {
	case model.Unity100:
		leadsCount = model.CollageSize100
	case model.Unity1K:
		leadsCount = model.CollageSize1K
	case model.Unity10K:
		leadsCount = model.CollageSize10K
	case model.Unity100K:
		leadsCount = model.CollageSize100K
	default:
		log.Fatal().Msgf("[frigg] НЕПОНЯТНЫЙ УРОВЕНЬ ЕДИНСТВА %s - %d", unity.Mask, unity.Rank)
	}

	state.Active().TotalLeads = uint(leadsCount)
	f.sendState(ctx, state)

	leads := make([]uint, 0, leadsCount)
	for i := 0; i < leadsCount; i++ {
		lead, _, err := f.muninn.RememberUnityLead(ctx, applicants)
		if err != nil {
			return model.Unity{},
				errors.Wrapf(err, "[frigg] ОШИБКА ВЫБОРА ЛИДЕРА %d/%d ДЛЯ ЕДИНСТВА %s", i+1, leadsCount, unity.Mask)
		}
		leads = append(leads, lead)
		state.Active().Leads = append(state.Active().Leads, lead)
		f.sendState(ctx, state)
	}

	log.Info().Msgf("[frigg] ДЛЯ ЕДИНСТВА %s ВЫБРАНЫ ЛИДЕРЫ %v", unity.Mask, leads)
	// Frigg: лидеры выбраны, нужно составить из картинок коллаж и затем сохранить его на warehouse
	// Frigg: само единство будет сохранено в БД последним

	maxArtID, err := f.artPile.GetMaxArtID(ctx)
	if err != nil {
		return model.Unity{}, errors.Wrapf(err, "[frigg] ОШИБКА ДОСТУПА К МАКСИМАЛЬНОМУ ID, %s", unity.Mask)
	}

	leadsB, err := json.Marshal(leads)
	if err != nil {
		return model.Unity{}, errors.Wrapf(err, "[frigg] ОШИБКА УПАКОВКИ ЛИДЕРОВ %s", unity.Mask)
	}
	unity.Leads = string(leadsB)
	unity.Version = unity.Version + 1

	state.Active().CollageStarted = true
	state.Active().Unity = &unity
	f.sendState(ctx, state)

	img, err := f.collage.MakeCollage(ctx, unity.Mask, leads, maxArtID)
	if err != nil {
		return model.Unity{}, errors.Wrapf(err, "[frigg] ОШИБКА СБОРКИ КОЛЛАЖА ДЛЯ %s", unity.Mask)
	}

	if err := f.collage.SaveCollage(ctx, unity, img); err != nil {
		return model.Unity{}, errors.Wrapf(err, "[frigg] ОШИБКА СОХРАНЕНИЯ КОЛЛАЖА ДЛЯ %s", unity.Mask)
	}

	state.Active().CollageFinished = true
	f.sendState(ctx, state)

	if maxArtID >= unity.MaxID {
		unity.State = model.UnityStateUnified
	} else {
		unity.State = model.UnityStatePreUnified
	}
	unity, err = f.unityPile.Save(ctx, unity)
	if err != nil {
		return model.Unity{}, errors.Wrapf(err, "[frigg] СБОЙ СОХРАНЕНИЯ ОБЪЕДИНЁННОГО ЕДИНСТВА %s", unity.Mask)
	}
	log.Info().Msgf("[frigg] ЕДИНСТВО %s ОБЪЕДИНЕНО. СТАТУС: %s. НОВАЯ ВЕРСИЯ: %d", unity.Mask, unity.State, unity.Version)

	state.Active().Unity = &unity
	f.sendState(ctx, state)

	s := time.Now()
	state.Active().CurrentEnjoyTime = 0
	state.Active().ExpectedEnjoyTime = f.unificationEnjoyTimeSec
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(time.Second):
				state.Active().CurrentEnjoyTime = uint(math.Ceil(time.Now().Sub(s).Seconds()))
				f.sendState(ctx, state)
				if state.Active().CurrentEnjoyTime >= state.Active().ExpectedEnjoyTime {
					return // цикл окончен
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msgf("[frigg] ПРЕЖДЕВРЕМЕННЫЙ ОСТАНОВ ПО КОНТЕКСТУ")
		return unity, nil
	case <-time.After(time.Second * time.Duration(f.unificationEnjoyTimeSec)):
	}

	state.ClearSubprocess()

	return unity, nil
}

// collectApplicants
// Frigg: Выбор всех возможных претендентов на лидерство в виде массива ID-номеров. Далее Один (Мунин) выберет конкретных
// Frigg: если это сотенное единство, то тут появится массив всех возможных ID [100, 101, 102,... 199]
func (f *Frigg) collectApplicants(ctx context.Context, unity model.Unity) ([]uint, error) {
	var applicants []uint
	if unity.Rank == model.Unity100 {
		for i := unity.MinID; i <= unity.MaxID; i++ {
			applicants = append(applicants, i)
		}
		return applicants, nil
	}

	children, err := f.unityPile.GetChildren(ctx, unity)
	if err != nil {
		return applicants, errors.Wrapf(err, "[frigg] ОШИБКА ВЫБОРА ДОЧЕРНИХ ЕДИНСТВ ДЛЯ %s", unity.Mask)
	}

	for _, child := range children {
		var leads []uint
		if err := json.Unmarshal([]byte(child.Leads), &leads); err != nil {
			return applicants, errors.Wrapf(err, "[frigg] ОШИБКА РАСШИФРОВКИ ЛИДЕРОВ ЕДИНСТВА %s", child.Mask)
		}
		applicants = append(applicants, leads...)
	}
	log.Info().Msgf("[frigg] ДЛЯ ЕДИНСТВА %s ВЫБРАНЫ ПРЕДЕНТЕНДЫ: %+v", unity.Mask, applicants)
	return applicants, nil
}

func (f *Frigg) sendState(ctx context.Context, state *model.FriggState) {
	if err := f.heimdallr.SendFriggState(ctx, *state); err != nil {
		log.Error().Err(err).Msgf("[frigg] ХЕЙМДАЛЛЬ НЕ МОЖЕТ ОТПРАВИТЬ ДРАККАР В АЛЬВХЕЙМ!")
	}
}
