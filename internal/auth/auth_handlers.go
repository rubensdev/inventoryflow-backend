package auth

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rubensdev/inventoryflow-backend/internal/jsonutil"
)

type AuthHandler struct {
	logger            *log.Logger
	authService       AuthService
	accessTokenSecret string
}

func NewAuthHandler(logger *log.Logger, authService AuthService, accessTokenSecret string) *AuthHandler {
	return &AuthHandler{
		logger:            logger,
		authService:       authService,
		accessTokenSecret: accessTokenSecret,
	}
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginReq := NewLoginRequest()

	jsonRes := jsonutil.NewJSONResponse(h.logger)

	err := jsonutil.ReadJSON(w, r, &loginReq)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	valid := loginReq.Validate()
	if !valid {
		err := jsonutil.WriteJSON(w, http.StatusBadRequest, jsonutil.H{
			"errors": loginReq.GetErrors(),
		}, nil)
		if err != nil {
			jsonRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	user, err := h.authService.ValidateCredentials(*loginReq)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	claims := UserClaims{
		ID:        strconv.Itoa(int(user.ID)),
		Username:  user.Username,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	accessToken, err := NewAccessToken(claims, h.accessTokenSecret)
	if err != nil {
		jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	cookie := NewAccessTokenCookie(accessToken, DefaultCookieMaxAge)
	http.SetCookie(w, cookie)

	err = jsonutil.WriteJSON(w, http.StatusOK, jsonutil.H{"user": user}, nil)
	if err != nil {
		jsonRes.ServerErrorResponse(w, r, err)
	}
}
