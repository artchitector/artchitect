package pantheon

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"math"
	"sync"
	"time"
)

// Huginn "думающий". Один из двух воронов Odin-а. Huginn - это мысль Одина.
// Huginn: я могу объяснить энтропию, видимую пустым глазом Odin-а в ткани мироздания.
// Huginn: я могу превратить её в мысль Odin-а, представленную в программе как числа (uint64, и затем float64).
// Huginn: набор таких чисел позволит Odin-у придумать идею картины, если он достаточно долго будет смотреть в ткань мироздания.
type Huginn struct {
	// Huginn тоже видит энтропию, которую и сам Odin видит в своём LostEye.
	// Huginn пробрасывает цепочку вызовов в глаз
	lostEye *LostEye

	// Huginn: Я отправляю энтропию не только Muninn, но и Heimdallr, который непрерывно ретранслирует её
	// 		по радужному мосту communication.Bifrost в Alfheimr (api-gateway).
	// Huginn: Далее из Alfheimr светлые эльфы переправят эту энтропию на Землю, в Midgard (frontend), где её увидят люди.
	// Huginn: для этого у меня тут механизм нескольких подписчиков
	// Odin: Воистину, пусть и смертные тоже увидят эту ткань пространства в виде меняющихся картинок.
	// Odin: Они всё равно ничего в ней не поймут.
	sMutex      sync.Mutex
	subscribers []*huSubscriber

	lastEntropy *model.EntropyPackExtended // может быть пустым, если кто-то уже воспользовался этой энтропией (она на один раз)
}

// Heimdallr подписывается на данные от Huginn и получает энтропию по подписке на go-канал
// huSubscriber - механизм подписок
type huSubscriber struct {
	ctx context.Context
	ch  chan model.EntropyPackExtended
}

func NewHuginn(lostEye *LostEye) *Huginn {
	return &Huginn{
		lostEye:     lostEye,
		sMutex:      sync.Mutex{},
		subscribers: make([]*huSubscriber, 0),
		lastEntropy: nil,
	}
}

func (h *Huginn) StartEntropyRealize(ctx context.Context) {
	lostEyeChan := h.lostEye.Subscribe(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("[huginn] ОСТАНАВЛИВАЮ СМОТРЕНИЕ В ПРОСТРАНСТВО")
			return
		case pack := <-lostEyeChan:
			// Huginn: у нас тут объект энтропии, но он не заполнен (там только картинки)
			// Huginn: я это преобразую в числа, но картинки тоже оставлю (их сохранят позже для истории)
			pack = h.realizeEntropy(ctx, pack)

			// Huginn: теперь посылку надо отправить подписчикам
			h.notifyListeners(ctx, pack) // Huginn: сообщаю всем заинтересованным, что новая энтропия осознана и посчитана

			// Huginn: но Muninn не требует постоянный поток энтропии. Для него я отложу последнюю энтропию, и он
			// воспользуется ей тогда, когда ему нужно будет дать Odin-у новое воспоминание.
			h.sMutex.Lock()
			h.lastEntropy = &pack
			h.sMutex.Unlock()
		}
	}
}

// Subscribe - отдаёт канал, из которого подписчик читает сообщения.
// Если подписчик закрывает контекст, то отправка прерывается.
func (h *Huginn) Subscribe(subscriberCtx context.Context) chan model.EntropyPackExtended {
	ch := make(chan model.EntropyPackExtended)
	sub := huSubscriber{
		ctx: subscriberCtx,
		ch:  ch,
	}
	h.sMutex.Lock()
	defer h.sMutex.Unlock()

	h.subscribers = append(h.subscribers, &sub)
	go func() {
		<-subscriberCtx.Done()
		h.unsubscribe(&sub)
	}()
	return ch
}

func (h *Huginn) unsubscribe(sub *huSubscriber) {
	h.sMutex.Lock()
	defer h.sMutex.Unlock()

	idx := slices.IndexFunc(h.subscribers, func(s *huSubscriber) bool { return s == sub })
	if idx == -1 {
		log.Warn().Msgf("[huginn] ПОДПИСАНТ ИСЧЕЗ. ПРОБЛЕМА")
		return
	}

	h.subscribers = append(h.subscribers[:idx], h.subscribers[idx+1:]...)
	log.Debug().Msgf("[huginn] ПОДПИСАНТ %d УДАЛЁН. УСПЕХ.", idx)
}

func (h *Huginn) notifyListeners(ctx context.Context, pack model.EntropyPackExtended) {
	h.sMutex.Lock()
	subscribers := h.subscribers[:]
	h.sMutex.Unlock()

	for _, sub := range subscribers {
		// отправка энтропии всем слушателям
		go func(s *huSubscriber) {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(time.Second):
				log.Error().Msgf("[huginn] ОТПРАВКА ЗАВИСЛА, ЭНТРОПИЯ ПОТЕРЯНА")
			case s.ch <- pack:
				//log.Debug().Msgf("[huginn] ЭНТРОПИЯ ОТПРАВЛЕНА")
			}
		}(sub)
	}
}

// realizeEntropy
// Huginn: рассматриваю две картинки энтропии (прямую и обратную)
// Huginn: и превращаю их в нерушимые и твёрдые сущности - конкретные числа, на которые Odin сможет положиться
func (h *Huginn) realizeEntropy(ctx context.Context, pack model.EntropyPackExtended) model.EntropyPackExtended {
	entropyVal := h.matrixToInt(pack.Entropy.Matrix)
	choiceVal := h.matrixToInt(pack.Choice.Matrix)

	pack.Entropy.IntValue = entropyVal
	pack.Entropy.FloatValue = float64(entropyVal) / float64(math.MaxUint64)
	pack.Choice.IntValue = choiceVal
	pack.Choice.FloatValue = float64(choiceVal) / float64(math.MaxUint64)

	return pack
}

// matrixToInt - ворон Huginn будет превращать картинку (матрицу сил пикселей) в число
// Huginn: на картине 64 пикселя, каждый светится с силой от 0 до 255.
// Huginn: Я превращу каждый пиксель в 0 или 1 (смотря насколько он светится, больше ли чем наполовину).
// Huginn: Эти 64 включенных или выключенных пикселя станут битами в uint64-числе
func (h *Huginn) matrixToInt(matrix model.EntropyMatrix) uint64 {
	var result uint64
	for x := 0; x < matrix.Size(); x++ {
		for y := 0; y < matrix.Size(); y++ {
			power := matrix.Get(x, y)
			isEnabledPixel := power >= math.MaxUint8/2
			if isEnabledPixel {
				byteIndex := x*matrix.Size() + y
				result = result | 1<<(63-byteIndex)
			}
		}
	}
	return result
}
