package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/sankethkini/NewsLetter-Backend/internal/enum"
)

type JWTConfig struct {
	Secret   string `yaml:"secret"`
	Duration int    `yaml:"duration"`
}

//go:generate mockgen -destination jwt_mock.go -package auth github.com/sankethkini/NewsLetter-Backend/pkg/auth JWTManager
type JWTManager interface {
	Generator(string, enum.Access) (string, error)
	Validate(string) (*UserClaims, error)
}

type jwtManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Email string
	Role  string
}

func NewJWTManager(cfg JWTConfig) JWTManager {
	return &jwtManager{
		secretKey:     cfg.Secret,
		tokenDuration: time.Hour * time.Duration(cfg.Duration),
	}
}

// generate new token.
func (manager *jwtManager) Generator(email string, role enum.Access) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Email: email,
		Role:  role.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// validate the token.
func (manager *jwtManager) Validate(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken,
		&UserClaims{},
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token claims")
			}
			return []byte(manager.secretKey), nil
		})
	if err != nil {
		return nil, errors.Wrap(err, "invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
