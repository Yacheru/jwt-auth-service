package jwt

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"jwt-auth-service/init/logger"
	"jwt-auth-service/pkg/constants"
	"math/rand"
	"time"
)

type TokenManager interface {
	NewAccessToken(ipAddr string, userId string, ttl time.Duration) (string, error)
	ParseAccessToken(accessToken string) (string, error)
	NewRefreshToken() string
	GenerateRefreshBcrypt(refreshToken string) (string, error)
	CheckBcryptMatch(hashedToken, token string) error
}

type Manager struct {
	salt string
}

func NewJWTManager(salt string) *Manager {
	return &Manager{salt: salt}
}

func (m *Manager) NewAccessToken(ipAddr string, userId string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp": time.Now().Add(ttl).Unix(),
		"ip":  ipAddr,
		"sub": userId,
	})

	return token.SignedString([]byte(m.salt))
}

func (m *Manager) ParseAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(m.salt), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) NewRefreshToken() string {
	b := make([]byte, 16)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	r.Read(b)

	return fmt.Sprintf("%x", md5.New().Sum(b))
}

func (m *Manager) CheckBcryptMatch(hashedToken, token string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
}

func (m *Manager) GenerateRefreshBcrypt(refreshToken string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.MinCost)
	if err != nil {
		logger.Error(err.Error(), constants.JWTCategory)

		return "", err
	}

	return string(b), nil
}
