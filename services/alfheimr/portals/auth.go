package portals

import (
	"crypto/hmac"
	"crypto/sha256"

	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const FakeLocalToken = "FAKE_LOCAL_TOKEN"
const FakeUserID = 999

type AuthPortal struct {
	allowFakeAuth  bool
	secret         string
	artchitectHost string
	botToken       string
}

func NewAuthPortal(
	allowFakeAuth bool,
	secret string,
	artchitectHost string,
	botToken string,
) *AuthPortal {
	return &AuthPortal{
		allowFakeAuth:  allowFakeAuth,
		secret:         secret,
		artchitectHost: artchitectHost,
		botToken:       botToken,
	}
}

func (ap *AuthPortal) HandleLogin(c *gin.Context) {
	log.Info().Msgf("[auth:portal] ПОПЫТКА ВХОДА")
	if err := ap.checkFromTelegram(c.Request.URL.Query()); err != nil {
		log.Error().Err(err).Msgf("[auth:portal] ЗАПРОС НЕ ОТ ТЕЛЕГРАМ")
		c.JSON(http.StatusUnauthorized, wrapError(errors.Errorf("[auth:portal] НЕКОРРЕКТНЫЕ ДАННЫЕ ВХОДА")))
		return
	}

	j, err := ap.generateJWT(c.Request.URL.Query())
	if err != nil {
		log.Error().Err(err).Msgf("[auth:portal] ОШИБКА ГЕНЕРАЦИИ JWT-ТОКЕНА")
		c.JSON(http.StatusInternalServerError, wrapError(errors.Errorf("[auth:portal] ОШИБКА ГЕНЕРАЦИИ JWT")))
		return
	}

	params := url.Values{}
	params.Add("token", j)
	params.Add("username", c.Request.URL.Query().Get("username"))
	params.Add("photo_url", c.Request.URL.Query().Get("photo_url"))
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?%s", ap.artchitectHost, params.Encode()))
}

func (ap *AuthPortal) HandleMe(c *gin.Context) {
	userId := ap.getUserIdFromContext(c)
	c.JSON(http.StatusOK, userId)
}

func (ap *AuthPortal) getUserIdFromContext(c *gin.Context) uint {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || authHeader == "null" {
		return 0
	}

	if authHeader == FakeLocalToken {
		if !ap.allowFakeAuth {
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
		return ap.secret, nil
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

func (ap *AuthPortal) checkFromTelegram(values url.Values) error {
	/*
		https://core.telegram.org/widgets/login
		Data-check-string is a concatenation of all received fields, sorted in alphabetical order, in the format
		key=<value> with a line feed character ('\n', 0x0A) used as separator – e.g.,
		'auth_date=<auth_date>\nfirst_name=<first_name>\nid=<id>\nusername=<username>'.

		data_check_string = ...
		secret_key = SHA256(<bot_token>)
		if (hex(HMAC_SHA256(data_check_string, secret_key)) == hash) {
		  // data is from Telegram
		}
	*/

	hash := values.Get("hash")
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for key, _ := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	stringPieces := make([]string, 0, len(values))
	for _, key := range keys {
		stringPieces = append(stringPieces, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	dataCheckString := strings.Join(stringPieces, "\n")
	secretKey := makeSha256(ap.botToken)
	encryptedDataCheckString := makeHmacSha256([]byte(dataCheckString), secretKey)
	if hash != hex.EncodeToString(encryptedDataCheckString) {
		return errors.Errorf("[auth:portal]")
	}
	return nil
}

func makeSha256(str string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hasher.Sum(nil)
}

func makeHmacSha256(data []byte, key []byte) []byte {
	sig := hmac.New(sha256.New, key)
	sig.Write(data)
	return sig.Sum(nil)
}

func (ap *AuthPortal) generateJWT(v url.Values) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = v.Get("id")
	claims["first_name"] = v.Get("first_name")
	claims["username"] = v.Get("username")
	claims["photo_url"] = v.Get("photo_url")
	claims["auth_date"] = v.Get("auth_date")

	log.Info().Msgf("%s", ap.secret)
	tokenStr, err := token.SignedString(ap.secret)
	if err != nil {
		return "", errors.Wrapf(err, "[login_handler] failed to sign JWT")
	}
	return tokenStr, nil
}
