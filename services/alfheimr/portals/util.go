package portals

import (
	"math/rand"

	"github.com/gin-gonic/gin"
)

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func wrapError(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func makeRadioConnectionID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
