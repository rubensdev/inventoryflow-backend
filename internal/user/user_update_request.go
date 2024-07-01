package user

import (
	"fmt"
)

type UserUpdateRequest struct {
	UserRequest
	Version   int `json:"version"`
	validated bool
}

func NewUserUpdateRequest(firstname, lastname, username, email string, version int) *UserUpdateRequest {
	return &UserUpdateRequest{
		UserRequest: UserRequest{
			Firstname: firstname,
			Lastname:  lastname,
			Username:  username,
			Email:     email,
		},
		Version: version,
	}
}

func (r *UserUpdateRequest) Validate() (valid bool) {
	r.ValidateCommonFields()

	valid = len(r.Errors) == 0
	r.validated = valid

	return
}

func (r *UserUpdateRequest) ToUser() (*User, error) {
	if !r.validated {
		return nil, fmt.Errorf("user update data hasn't been validated")
	}

	return &User{
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Email:     r.Email,
		Username:  r.Username,
		Version:   r.Version,
	}, nil
}
