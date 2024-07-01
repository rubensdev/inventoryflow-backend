package user

import (
	"errors"
	"fmt"
)

var ErrEditConflict = errors.New("edit conflict")
var ErrUserNotFound = errors.New("user not found")

func DuplicatedEmailError(email string) error {
	return fmt.Errorf("the email %s is being used", email)
}

func DuplicatedUsernameError(username string) error {
	return fmt.Errorf("the username %s is being used", username)
}
