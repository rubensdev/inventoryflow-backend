package auth_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/rubensdev/inventoryflow-backend/internal/auth"
)

func TestRequireAuthenticationMiddleware(t *testing.T) {

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK!"))
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	accessTokenSecret := "TestAccessTokenSecretDoNotShare!"

	accessToken, err := auth.NewAccessToken(auth.UserClaims{
		ID:        "1",
		Username:  "foobar",
		Email:     "foobar@test.es",
		Firstname: "Foo",
		Lastname:  "Bar",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	},
		accessTokenSecret,
	)
	if err != nil {
		t.Fatalf("error creating new access token. %v", err)
	}

	middleware := auth.NewAuthMiddleware(accessTokenSecret, logger)

	endpointURL := "/v1/test"

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, endpointURL, middleware.RequireAuthenticatedUser(testHandler))

	req := httptest.NewRequest(http.MethodGet, endpointURL, nil)
	w := httptest.NewRecorder()

	accessTokenCookie := auth.NewAccessTokenCookie(accessToken, auth.DefaultCookieMaxAge)
	req.AddCookie(accessTokenCookie)

	router.ServeHTTP(w, req)

	t.Log(w.Code)
	t.Log(w.Body.String())
}
