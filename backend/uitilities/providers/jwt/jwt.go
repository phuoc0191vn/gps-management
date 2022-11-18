package jwt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

type JWTService struct {
	secretKey string
}

type JWTClaims struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

// IssueJWT generates a JWTService
func (service JWTService) IssueJWT(data interface{}) (string, error) {
	claims := JWTClaims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        bson.NewObjectId().Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(service.secretKey))
}

// IssueJWT generates a JWTService with expiration
func (service JWTService) IssueJWTWithExpiration(data interface{}, exp time.Time) (string, error) {
	claims := JWTClaims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        bson.NewObjectId().Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(service.secretKey))
}

// ValidateJWT validates a JWTService
func (service JWTService) ValidateJWT(tokenString string, out interface{}) error {
	token, e := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(service.secretKey), nil
	})
	if e != nil {
		return e
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return service.parseDataFromMap(claims, out)
	}

	return fmt.Errorf("invalid JWTService")
}

func (service JWTService) parseDataFromMap(m jwt.MapClaims, out interface{}) error {
	data, ok := m["data"]
	if !ok || data == nil {
		return fmt.Errorf("invalid JWTService Claims: Data")
	}

	buffers, e := json.Marshal(data)
	if e != nil {
		return e
	}

	return json.Unmarshal(buffers, out)
}

func NewJwtService(secretKey string) *JWTService {
	service := new(JWTService)
	service.secretKey = secretKey
	return service
}
