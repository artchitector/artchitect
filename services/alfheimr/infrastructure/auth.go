package infrastructure

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

const FakeLocalToken = "FAKE_LOCAL_TOKEN"
const FakeUserID = 999

type AuthService struct {
	allowFakeAuth bool
	secret        []byte
}

func NewAuthService(allowFakeAuth bool, secret []byte) *AuthService {
	return &AuthService{allowFakeAuth: allowFakeAuth, secret: secret}
}

func (as *AuthService) GetUserIdFromContext(c *gin.Context) uint {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || authHeader == "null" {
		return 0
	}

	if authHeader == FakeLocalToken {
		if !as.allowFakeAuth {
			err := errors.Errorf("[auth:portal] ФЕЙКОВАЯ АВТОРИЗАЦИЯ ЗАПРЕЩЕНА")
			log.Error().Err(err).Send()
			return 0
		}
		return FakeUserID
	}

	log.Debug().Msgf("header=%s", authHeader)

	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("[auth:portal] НЕОЖИДАННЫЙ МЕТОД ШИФРОВАНИЯ КЛЮЧА: %v", token.Header["alg"])
		}
		return as.secret, nil
	})

	if err != nil {
		log.Error().Err(err).Msgf("[auth:portal] ОШИБКА ЧТЕНИЯ ТОКЕНА")
		return 0
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ID, err := strconv.ParseInt(fmt.Sprintf("%s", claims["id"]), 10, 64)
		if err != nil {
			log.Error().Err(err).Msgf("[auth:portal] ОШИБКА ЧТЕНИЯ ID")
			return 0
		}
		log.Info().Msgf("[auth:portal] ВХОД ПОДТВЕРЖДЁН. ID=%d", ID)
		return uint(ID)
	} else {
		log.Error().Msgf("[auth:portal] ПРОБЛЕМА С РАСШИФРОВКОЙ ДАННЫХ КЛЮЧА")
		return 0
	}
}
