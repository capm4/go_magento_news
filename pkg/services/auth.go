package services

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserRole string
	NickName string
}

const (
	tokenTTL = 24 * time.Hour
)

func CreateJWT(secret string, userRole, nickName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, userRole, nickName,
	})
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		logrus.Warn(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func GetUserNickNameFromJWT(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*TokenClaims)

	return claims.NickName
}
