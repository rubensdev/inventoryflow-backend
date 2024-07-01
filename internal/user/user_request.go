package user

import (
	"regexp"
	"strings"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type UserRequest struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Errors    map[string]string
}

func (r *UserRequest) Sanitize() {
	r.Firstname = strings.TrimSpace(r.Firstname)
	r.Lastname = strings.TrimSpace(r.Lastname)
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
}

func (r *UserRequest) ValidateCommonFields() {
	r.Sanitize()

	if r.Firstname == "" {
		r.AddError("first_name", "firstname is required")
	}

	if r.Lastname == "" {
		r.AddError("last_name", "lastname is required")
	}

	if r.Username == "" {
		r.AddError("username", "username is required")
	}

	if r.Email == "" {
		r.AddError("email", "email is required")
	} else if !EmailRX.MatchString(r.Email) {
		r.AddError("email", "email is not valid")
	}
}

func (r *UserRequest) GetErrors() map[string]string {
	return r.Errors
}

func (r *UserRequest) AddError(key string, value string) {
	if r.Errors == nil {
		r.Errors = make(map[string]string)
	}
	r.Errors[key] = value
}
