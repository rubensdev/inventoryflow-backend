package auth

import "strings"

type LoginRequest struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Errors   map[string]string `json:"errors"`
}

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}

func (r *LoginRequest) Sanitize() {
	r.Username = strings.TrimSpace(r.Username)
}

func (r *LoginRequest) Validate() (valid bool) {
	r.Sanitize()

	if r.Username == "" {
		r.AddError("username", "The username is required")
	}
	if r.Password == "" {
		r.AddError("password", "The password is required")
	}

	valid = len(r.Errors) == 0
	return
}

func (r *LoginRequest) GetErrors() map[string]string {
	return r.Errors
}

func (r *LoginRequest) AddError(field string, msg string) {
	if r.Errors == nil {
		r.Errors = make(map[string]string)
	}
	r.Errors[field] = msg
}
