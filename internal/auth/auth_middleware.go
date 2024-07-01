package auth

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/rubensdev/inventoryflow-backend/internal/jsonutil"
	userdom "github.com/rubensdev/inventoryflow-backend/internal/user"
)

type AuthMiddleware struct {
	accessTokenSecret string
	logger            *log.Logger
}

func NewAuthMiddleware(accessToken string, logger *log.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		accessTokenSecret: accessToken,
		logger:            logger,
	}
}

func (m AuthMiddleware) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		jsonRes := jsonutil.NewJSONResponse(m.logger)

		var accessToken string
		for _, c := range cookies {
			if c.Name == AccessTokenCookieName {
				accessToken = c.Value
				break
			}
		}

		if accessToken == "" {
			jsonRes.ErrorResponse(w, r, http.StatusUnauthorized, "missing access token")
			return
		}

		claims, err := ParseAccessToken(accessToken, m.accessTokenSecret)
		if err != nil {
			jsonRes.ErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		userID, err := strconv.Atoi(claims.ID)
		if err != nil {
			jsonRes.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		user := &userdom.User{
			Username:  claims.Username,
			Firstname: claims.Firstname,
			Lastname:  claims.Lastname,
			Email:     claims.Email,
			ID:        int64(userID),
		}

		var userCtxKey UserCtxKey = "user"
		r = r.WithContext(context.WithValue(r.Context(), userCtxKey, user))

		next.ServeHTTP(w, r)
	})
}
