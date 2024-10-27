package auth

import "time"

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

}

