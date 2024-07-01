package user

import (
	"fmt"
	"sync"
)

type InMemoryUserRepository struct {
	users []*User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make([]*User, 0),
	}
}

func (r *InMemoryUserRepository) GetAll() ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.users, nil
}

func (r *InMemoryUserRepository) GetByID(id int64) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.users {
		if r.users[i].ID == id {
			return r.users[i], nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) GetByEmail(email string) *User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.users {
		if r.users[i].Email == email {
			return r.users[i]
		}
	}
	return nil
}

func (r *InMemoryUserRepository) GetByUsername(username string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.users {
		if r.users[i].Username == username {
			return r.users[i], nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) Create(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = int64(len(r.users)) + 1

	for _, u := range r.users {
		if u.Username == user.Username {
			return fmt.Errorf("an user with username \"%s\" already exists", user.Username)
		}
		if u.Email == user.Email {
			return fmt.Errorf("an user with email \"%s\" already exists", user.Email)
		}
	}

	r.users = append(r.users, user)
	return nil
}

func (r *InMemoryUserRepository) Update(user *User) error {
	u := r.GetByEmail(user.Email)
	if u != nil && u.ID != user.ID {
		return DuplicatedEmailError(user.Email)
	}

	u, err := r.GetByUsername(user.Username)
	if err != nil {
		return err
	}
	if u != nil && u.ID != user.ID {
		return DuplicatedUsernameError(user.Username)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.users {
		u := r.users[i]
		if u.ID == user.ID {
			if u.Version != user.Version {
				return ErrEditConflict
			}
			user.Version++
			u.Username = user.Username
			u.Email = user.Email
			u.Firstname = user.Firstname
			u.Lastname = user.Lastname
			u.Version = user.Version
			return nil
		}
	}
	return nil
}

func (r *InMemoryUserRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	filteredUsers := make([]*User, 0, len(r.users)-1)

	found := false

	for i := range r.users {
		if r.users[i].ID != id {
			filteredUsers = append(filteredUsers, r.users[i])
			continue
		}
		found = true
	}

	if !found {
		return ErrUserNotFound
	}

	r.users = filteredUsers
	return nil
}
