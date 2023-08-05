package pantheon

import (
	"context"
	"github.com/artchitector/artchitect2/model"
	"github.com/pkg/errors"
)

// Muninn - "помнящий". Один из воронов Odin-а.
// Muninn: я память Odin-а Всеотца, я могу вспомнить всё, что ему нужно для придумывания идеи картины.
// Mininn: я буду использовать интерпретацию энтропии от Huginn, и дам Odin-у нужные ему воспоминания
// Odin: я существую вне времени, и заранее знаю, что и когда будет нарисовано. Мне просто нужно вспомнить, что именно в этот раз.
type Muninn struct {
	// Muninn спрашивает Huginn о том, что Odin видит в своём LostEye
	// Muninn пробрасывает цепочку вызовов в Huginn, а оттуда в глаз Odin-а - LostEye
	huginn *Huginn
}

func (m *Muninn) rememberSeed(ctx context.Context) (uint, model.EntropyPack, error) {
	return 0, model.EntropyPack{}, errors.Errorf("[muninn] НЕ РЕАЛИЗОВАНО ЕЩЕ")
}

func (m *Muninn) rememberNumberOfWords(ctx context.Context) (uint, model.EntropyPack, error) {
	return 0, model.EntropyPack{}, errors.Errorf("[muninn] НЕ РЕАЛИЗОВАНО ЕЩЕ")
}

// rememberWord - вспомнить слово. model.Word уже содержит в себе model.EntropyPack внутри, потому он и не возвращается
func (m *Muninn) rememberWord(ctx context.Context) (model.Word, error) {
	return model.Word{}, errors.Errorf("[muninn] НЕ РЕАЛИЗОВАНО ЕЩЕ")
}
