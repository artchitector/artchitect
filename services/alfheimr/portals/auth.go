package portals

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

const FakeLocalToken = "FAKE_LOCAL_TOKEN"
const FakeUserID = 999

type AuthPortal struct {
	allowFakeAuth bool
	secret        string
}

func NewAuthPortal(allowFakeAuth bool, secret string) *AuthPortal {
	return &AuthPortal{allowFakeAuth: allowFakeAuth, secret: secret}
}

func (ah *AuthPortal) HandleMe(c *gin.Context) {
	userId := ah.getUserIdFromContext(c)
	c.JSON(http.StatusOK, userId)
}

func (ah *AuthPortal) getUserIdFromContext(c *gin.Context) uint {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || authHeader == "null" {
		return 0
	}

	if authHeader == FakeLocalToken {
		if !ah.allowFakeAuth {
			err := errors.Errorf("[auth] ФЕЙКОВАЯ АВТОРИЗАЦИЯ ЗАПРЕЩЕНА")
			log.Error().Err(err).Send()
			return 0
		}
		return FakeUserID
	}

	log.Debug().Msgf("header=%s", authHeader)

	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("[auth] НЕОЖИДАННЫЙ МЕТОД ШИФРОВАНИЯ КЛЮЧА: %v", token.Header["alg"])
		}
		return ah.secret, nil
	})

	if err != nil {
		log.Error().Err(err).Msgf("[auth] ОШИБКА ЧТЕНИЯ ТОКЕНА")
		return 0
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ID, err := strconv.ParseInt(fmt.Sprintf("%s", claims["id"]), 10, 64)
		if err != nil {
			log.Error().Err(err).Msgf("[auth] ОШИБКА ЧТЕНИЯ ID")
			return 0
		}
		log.Info().Msgf("[auth] ВХОД ПОДТВЕРЖДЁН. ID=%d", ID)
		return uint(ID)
	} else {
		log.Error().Msgf("[auth] ПРОБЛЕМА С РАСШИФРОВКОЙ ДАННЫХ КЛЮЧА")
		return 0
	}
}
