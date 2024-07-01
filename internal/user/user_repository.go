package user

type UserRepository interface {
	GetByID(id int64) (*User, error)
	GetByUsername(username string) (*User, error)
	GetAll() ([]*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int64) error
}
