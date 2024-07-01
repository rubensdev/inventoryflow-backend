package user

import (
	"unicode/utf8"
)

type UserPasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	Errors          map[string]string
}

func (r *UserPasswordRequest) Validate() (valid bool) {
	if r.CurrentPassword == "" {
		r.AddError("current_password", "current password is required")
	} else if utf8.RuneCountInString(r.CurrentPassword) < 8 {
		r.AddError("current_password", "current password must contain at least 8 characters")
	}

	if r.NewPassword == "" {
		r.AddError("password", "new password is required")
	} else if utf8.RuneCountInString(r.NewPassword) < 8 {
		r.AddError("password", "current password must contain at least 8 characters")
	}

	valid = len(r.Errors) == 0
	return
}

func (r *UserPasswordRequest) AddError(field string, value string) {
	if r.Errors == nil {
		r.Errors = make(map[string]string)
	}
	r.Errors[field] = value
}
