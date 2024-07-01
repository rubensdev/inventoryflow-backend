package user

import (
	"fmt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s UserService) GetAll() ([]*User, error) {
	return s.repo.GetAll()
}

func (s UserService) GetByID(id int64) (*User, error) {
	return s.repo.GetByID(id)
}

func (s UserService) GetByUsername(username string) (*User, error) {
	return s.repo.GetByUsername(username)
}

func (s UserService) Create(user *User) error {
	return s.repo.Create(user)
}

func (s UserService) Update(id int64, userReq *UserUpdateRequest) (*User, error) {
	_, err := s.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching user for update, %w", err)
	}

	user, err := userReq.ToUser()
	if err != nil {
		return nil, fmt.Errorf("error creating user model for update, %w", err)
	}

	// This is the part where we keep the data that we don't want to update. In this case
	// we just pass the ID to the validated data from the UserUpdateRequest.
	user.ID = id

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s UserService) Delete(id int64) error {
	return s.repo.Delete(id)
}
