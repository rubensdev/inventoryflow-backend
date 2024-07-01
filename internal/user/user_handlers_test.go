package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	userdom "github.com/rubensdev/inventoryflow-backend/internal/user"
)

func TestA_User_Can_Be_Registered(t *testing.T) {
	firstname := "Foo"
	lastname := "Bar"
	username := "foobar"
	email := "foo@bar.baz"
	plainPasswd := "foobar$1234"

	// We assumed the user registration request was tested, and therefore, its validation is working.
	var registerForm userdom.UserRegisterRequest = *userdom.NewUserRegisterRequest(
		firstname,
		lastname,
		username,
		email,
		plainPasswd,
		plainPasswd,
	)

	data, err := json.Marshal(registerForm)
	if err != nil {
		t.Fatalf("error marshaling json %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(data))
	w := httptest.NewRecorder()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	inMemUserRepo := userdom.NewInMemoryUserRepository()
	userHandler := userdom.NewUserHandler(logger, *userdom.NewUserService(inMemUserRepo))

	userHandler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response struct {
		Message string
		User    *userdom.User
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshaling response json, %v", err)
	}

	if response.Message == "" {
		t.Errorf("a success message is expected")
	} else if response.Message != userdom.UserCreatedMsg {
		t.Errorf("expected %s, got %s", userdom.UserCreatedMsg, response.Message)
	}

	user := response.User

	if user == nil {
		t.Fatal("registed user data is expected")
	}

	if user.Firstname != firstname {
		t.Errorf("expected %s, got %s", firstname, user.Firstname)
	}

	if user.Lastname != lastname {
		t.Errorf("expected %s, got %s", lastname, user.Lastname)
	}

	if user.Username != username {
		t.Errorf("expected %s, got %s", username, user.Username)
	}

	if user.Email != email {
		t.Errorf("expected %s, got %s", email, user.Email)
	}

}

func TestFindUserByID(t *testing.T) {
	plaintextPasswd := "plaintext"

	newUser := &userdom.User{
		Firstname: "Foo",
		Lastname:  "Bar",
		Email:     "foo@bar.baz",
	}
	newUser.Password.Set(plaintextPasswd)

	inMemUserRepo := userdom.NewInMemoryUserRepository()

	err := inMemUserRepo.Create(newUser)
	if err != nil {
		t.Fatalf("error registering user %v", err)
	}

	userHandler := userdom.NewUserHandler(nil, *userdom.NewUserService(inMemUserRepo))

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", userHandler.GetUserByID)

	endpoint := fmt.Sprintf("/v1/users/%d", newUser.ID)

	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		User *userdom.User
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshaling response json, %v", err)
	}

	if response.User.Email != newUser.Email {
		t.Errorf("expected user email %s, got %s", newUser.Email, response.User.Email)
	}
}

func TestUserInfoCanBeUpdated(t *testing.T) {
	users := []*userdom.User{
		{
			Firstname: "Foo",
			Lastname:  "Bar",
			Username:  "foobar",
			Email:     "foo@bar.baz",
		},
		{
			Firstname: "Qux",
			Lastname:  "Toto",
			Username:  "quxtoto",
			Email:     "quxtoto@test.es",
		},
	}
	updateData := userdom.NewUserUpdateRequest(
		"UpdatedFoo",
		"UpdatedBar",
		"updatedfoobar",
		"foobar@updated.es",
		0,
	)

	inMemUserRepo := userdom.NewInMemoryUserRepository()

	// Register the users
	for i := range users {
		err := inMemUserRepo.Create(users[i])
		if err != nil {
			t.Fatalf("error creating a new user %v", err)
		}
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	userHandler := userdom.NewUserHandler(logger, *userdom.NewUserService(inMemUserRepo))

	router := httprouter.New()
	router.HandlerFunc(http.MethodPut, "/v1/users/:id", userHandler.UpdateUserByID)

	tests := []struct {
		Name               string
		Data               userdom.UserUpdateRequest
		Endpoint           string
		ExpectedErrorMsg   string
		ExpectedStatusCode int
	}{
		{
			Name: "update with existing email should return an error",
			Data: *userdom.NewUserUpdateRequest(
				"Baz",
				"Bae",
				users[0].Username,
				users[1].Email,
				0,
			),
			Endpoint:           fmt.Sprintf("/v1/users/%d", users[0].ID),
			ExpectedErrorMsg:   userdom.DuplicatedEmailError(users[1].Email).Error(),
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name: "update with existing username should return an error",
			Data: *userdom.NewUserUpdateRequest(
				"Chota",
				"Desu",
				users[1].Username,
				users[0].Email,
				0,
			),
			Endpoint:           fmt.Sprintf("/v1/users/%d", users[0].ID),
			ExpectedErrorMsg:   userdom.DuplicatedUsernameError(users[1].Username).Error(),
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name: "update user data with a non existing username and email should be successful",
			Data: *userdom.NewUserUpdateRequest(
				updateData.Firstname,
				updateData.Lastname,
				updateData.Username,
				updateData.Email,
				0,
			),
			Endpoint:           fmt.Sprintf("/v1/users/%d", users[0].ID),
			ExpectedErrorMsg:   "",
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name: "updating an user does not exist should return a user not found error message",
			Data: *userdom.NewUserUpdateRequest(
				updateData.Firstname,
				updateData.Lastname,
				updateData.Username,
				updateData.Email,
				0,
			),
			Endpoint:           "/v1/users/111",
			ExpectedErrorMsg:   userdom.ErrUserNotFound.Error(),
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name: "an edit conflict error should happen if I try to update an updated user",
			Data: *userdom.NewUserUpdateRequest(
				updateData.Firstname,
				updateData.Lastname,
				updateData.Username,
				updateData.Email,
				0,
			),
			Endpoint:           fmt.Sprintf("/v1/users/%d", users[0].ID),
			ExpectedErrorMsg:   userdom.ErrUserNotFound.Error(),
			ExpectedStatusCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := json.Marshal(tt.Data)
			if err != nil {
				t.Fatalf("error marshaling json data %v", err)
			}
			req := httptest.NewRequest(http.MethodPut, tt.Endpoint, bytes.NewReader(data))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			var response struct {
				Error string
				User  *userdom.User
			}

			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("error unmarshaling response JSON: %v", err)
			}

			if w.Code != tt.ExpectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.ExpectedStatusCode, w.Code)
			}

			switch w.Code {
			case http.StatusBadRequest:
				if response.Error != tt.ExpectedErrorMsg {
					t.Errorf("expected error message \"%s\", got \"%s\" ", tt.ExpectedErrorMsg, response.Error)
				}
			case http.StatusOK:
				assertEqual(t, updateData.Email, response.User.Email)
				assertEqual(t, updateData.Firstname, response.User.Firstname)
				assertEqual(t, updateData.Lastname, response.User.Lastname)
				assertEqual(t, updateData.Username, response.User.Username)
			}
		})
	}

}

func TestUserCanBeDeleted(t *testing.T) {
	// Create two users
	users := []*userdom.User{
		{
			Firstname: "User1FN",
			Lastname:  "User1LN",
			Username:  "user1",
			Email:     "user1@test.es",
		},
		{
			Firstname: "User2FN",
			Lastname:  "User2LN",
			Username:  "user2",
			Email:     "user2@test.es",
		},
	}

	inMemUserRepo := userdom.NewInMemoryUserRepository()
	for i := range users {
		inMemUserRepo.Create(users[i])
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	userHandler := userdom.NewUserHandler(logger, *userdom.NewUserService(inMemUserRepo))

	router := httprouter.New()
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", userHandler.DeleteUserByID)

	endpoint := fmt.Sprintf("/v1/users/%d", users[0].ID)
	req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", w.Code, http.StatusOK)
	}

	var response struct {
		Message string `json:"message"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshaling JSON %v", err)
	}

	assertEqual(t, userdom.UserDeletedMsg, response.Message)
}

func assertEqual(t *testing.T, expected, got string) {
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
