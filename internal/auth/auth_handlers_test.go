package auth_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rubensdev/inventoryflow-backend/internal/auth"
	userdom "github.com/rubensdev/inventoryflow-backend/internal/user"
)

func TestAuthenticateUser(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Create a new user and insert it into the in memory database
	userPasswd := "foobar$123"
	user := &userdom.User{
		Firstname: "Foo",
		Lastname:  "Bar",
		Username:  "foobar",
		Email:     "foobar@test.es",
	}
	user.Password.Set(userPasswd)

	inMemUserRepo := userdom.NewInMemoryUserRepository()
	inMemUserRepo.Create(user)

	// Setup the user service and the auth service
	userSrv := userdom.NewUserService(inMemUserRepo)
	authSrv := auth.NewAuthService(*userSrv)

	accessTokenSecret := "HelloTest123.ThisIsATest!"

	// Setup the Auth handlers (login)
	authHandler := auth.NewAuthHandler(logger, *authSrv, accessTokenSecret)

	// Marshaling the credentials
	loginReq := auth.LoginRequest{
		Username: user.Username,
		Password: userPasswd,
	}

	credentials, err := json.Marshal(loginReq)
	if err != nil {
		t.Fatalf("error marshaling JSON. %v", err)
	}

	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/v1/login", authHandler.Login)

	req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(credentials))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		User userdom.User `json:"user"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshaling response JSON, %v", err)
	}

	t.Logf("User is returned %+v", response.User)

	var accessToken string

	for _, c := range w.Result().Cookies() {
		if c.Name == auth.AccessTokenCookieName {
			accessToken = c.Value
			break
		}
	}

	userClaims, err := auth.ParseAccessToken(accessToken, accessTokenSecret)
	if err != nil {
		t.Fatalf("error parsing access token. %v", err)
	}

	if userClaims.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, userClaims.Email)
	}

	if userClaims.Firstname != user.Firstname {
		t.Errorf("expected firstname %s, got %s", user.Firstname, userClaims.Firstname)
	}

	if userClaims.Lastname != user.Lastname {
		t.Errorf("expected lastname %s, got %s", user.Lastname, userClaims.Lastname)
	}

	if userClaims.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, userClaims.Username)
	}
}
