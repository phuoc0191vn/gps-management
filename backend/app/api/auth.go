package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/model"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

type MyCustomClaims struct {
	jwt.StandardClaims
	Scopes    []string `json:"scopes"`
	DeviceIDs []string `json:"deviceIDs"`
}

type BasicAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *BasicAuthRequest) Trim() {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

type AuthHandler struct {
	JwtSecret string

	AccountRepository  repository.AccountRepository
	AccessIDRepository repository.UserAccessIDRepository
}

func (h *AuthHandler) BasicLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data BasicAuthRequest
	if err := BindJSON(r, &data); err != nil {
		ResponseError(w, r, err)
		return
	}
	data.Trim()

	if data.Email == "" || data.Password == "" {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "email or password is empty",
			Code:    http.StatusBadRequest,
		})
		return
	}

	account, err := h.AccountRepository.FindByEmail(data.Email)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: "account is not existed",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(data.Password))
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, ResponseBody{
			Message: "email or password is incorrect",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	claims := h.makeClaims(*account)
	tokenString, err := h.createToken(claims)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: "can not create token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if err := h.AccessIDRepository.Set(account.Email, claims.Id, fmt.Sprintf("%d", claims.ExpiresAt)); err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: "Unable to create token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Successfully login",
		Data: map[string]interface{}{
			"token": tokenString,
		},
	})
}

func (h *AuthHandler) makeClaims(account model.Account) MyCustomClaims {
	return MyCustomClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        bson.NewObjectId().Hex(),
			Audience:  account.UserID,
			Subject:   account.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(10, 0, 0).Unix(),
		},
		Scopes:    account.Scopes,
		DeviceIDs: account.DeviceIDs,
	}
}

func (h *AuthHandler) createToken(claims MyCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(h.JwtSecret))
}
