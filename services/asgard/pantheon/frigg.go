package pantheon

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var unityRanks = []uint{model.Unity100, model.Unity1K, model.Unity10K, model.Unity100K}
var updateIntervals = []uint{model.UpdateInterval100, model.UpdateInterval1K, model.UpdateInterval10K, model.UpdateInterval100K}

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
	unityPile unityPile
}

func NewFrigg(unityPile unityPile) *Frigg {
	return &Frigg{unityPile: unityPile}
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
	return false, errors.Errorf("[frigg] ОБЪЕДИНЕНИЕ ЕЩЕ НЕ ГОТОВО")
}

func (f *Frigg) createNonexistentUnities(ctx context.Context, art model.Art) error {
	for _, rank := range unityRanks {
		mask := art.GetUnityMask(rank)
		_, err := f.unityPile.Get(ctx, mask)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(err, "[frigg] СБОЙ ПОИСКА ЕДИНСТВА %s", mask)
		} else if err == nil {
			log.Debug().Msgf("[frigg] ЕДИНСТВО %s УЖЕ СОЗДАНО. ПРОПУСК", mask)
			continue
		}
		min, max, err := f.getMinMaxID(mask)
		if err != nil {
			return errors.Wrapf(err, "[frigg] МАТЕМАТИЧЕСКИЙ СБОЙ MIN-MAX %s", mask)
		}
		unity, err := f.unityPile.Create(ctx, mask, rank, min, max)
		if err != nil {
			return errors.Wrapf(err, "[frigg] СБОЙ СОЗДАНИЯ ЕДИНСТВА %s", mask)
		}
		log.Debug().Msgf("[frigg] СОЗДАНО ЕДИНСТВО %s", unity.Mask)
	}

	return nil
}

func (f *Frigg) getMinMaxID(mask string) (min uint, max uint, err error) {
	var min64, max64 uint64
	min64, err = strconv.ParseUint(
		strings.ReplaceAll(strings.ReplaceAll(mask, "X", "0"), "U", ""),
		0,
		64,
	)
	if err != nil {
		return 0, 0, err
	}
	max64, err = strconv.ParseUint(
		strings.ReplaceAll(strings.ReplaceAll(mask, "X", "9"), "U", ""),
		0,
		64,
	)
	return uint(min64), uint(max64), nil
}
