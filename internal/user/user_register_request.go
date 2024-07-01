package user

import "unicode/utf8"

type UserRegisterRequest struct {
	UserRequest
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	validated       bool
}

func NewUserRegisterRequest(firstname, lastname, username, email, password, passwordConfirm string) *UserRegisterRequest {
	return &UserRegisterRequest{
		UserRequest: UserRequest{
			Firstname: firstname,
			Lastname:  lastname,
			Username:  username,
			Email:     email,
		},
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

func (r *UserRegisterRequest) Validate() (valid bool) {
	r.ValidateCommonFields()

	if r.Password == "" {
		r.AddError("password", "password is required")
	} else if utf8.RuneCountInString(r.Password) < 8 {
		r.AddError("password", "password must have at least 8 characters")
	}

	if r.PasswordConfirm == "" {
		r.AddError("password_confirm", "password confirmation is required")
	} else if r.Password != r.PasswordConfirm {
		r.AddError("password_confirm", "passwords mismatch")
	}

	valid = len(r.Errors) == 0
	r.validated = valid

	return
}

func (r *UserRegisterRequest) GetModel() *User {
	user := &User{
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Email:     r.Email,
		Username:  r.Username,
	}

	user.Password.Set(r.Password)
	return user
}
