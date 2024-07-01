package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Cost parameter for generating the password hash. The higher the cost, the
// slower and more computationally expensive it is to generate the hash.
const costParam = 12

type User struct {
	ID        int64    `json:"id"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Version   int      `json:"version"`
	Password  password `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPasswd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPasswd), costParam)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPasswd
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPasswd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPasswd))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
