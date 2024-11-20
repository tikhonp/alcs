package auth

import (
	"time"

	"github.com/markbates/goth"
)

type User struct {
	Id          int       `db:"id"`
	Password    string    `db:"password"`
	LastLogin   time.Time `db:"last_login"`
	IsSuperuser bool      `db:"is_superuser"`
	Email       string    `db:"email"`
	PhoneNumber string    `db:"phone_number"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	IsActive    bool      `db:"is_active"`
	DateJoined  time.Time `db:"date_joined"`
}

// User Actions
type Users interface {

	// GetById returns a user by id
	GetById(id int) (*User, error)

	// IsUserHasPermissions checks if user has specified permissions
	IsUserHasPermissions(userId int, permissionCodenames ...string) (bool, error)

	// FromOAuth finds user by goth user or creates new one
	FromOAuth(guser *goth.User) (*User, error)
}

type users struct {

}

func NewUsers() Users {
    return &users{}
}

func (u *users) GetById(id int) (*User, error) {
    panic("Not implemented")
}

func (u *users) IsUserHasPermissions(userId int, permissionCodenames ...string) (bool, error) {
    panic("Not implemented")
}

func (u *users) FromOAuth(guser *goth.User) (*User, error) {
    panic("Not implemented")
}
