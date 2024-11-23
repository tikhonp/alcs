package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
)

type User struct {
	Id          int            `db:"id"`
	Password    sql.NullString `db:"password"`
	LastLogin   sql.NullTime   `db:"last_login"`
	IsSuperuser bool           `db:"is_superuser"`
	Email       string         `db:"email"`
	PhoneNumber sql.NullString `db:"phone_number"`
	FirstName   sql.NullString `db:"first_name"`
	LastName    sql.NullString `db:"last_name"`
	IsActive    bool           `db:"is_active"`
	DateJoined  time.Time      `db:"date_joined"`
	OAuthId     sql.NullString `db:"oauth_id"`
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
	db sqlx.DB
}

func NewUsers(db sqlx.DB) Users {
	return &users{db: db}
}

func (u *users) GetById(id int) (*User, error) {
	panic("Not implemented")
}

func (u *users) IsUserHasPermissions(userId int, permissionCodenames ...string) (bool, error) {
	panic("Not implemented")
}

func (u *users) FromOAuth(guser *goth.User) (*User, error) {

	var user User
	err := u.db.Get(&user, `SELECT (id, password, last_login, is_superuser, email, phone_number, first_name, last_name, is_active, date_joined, oauth_id) FROM users WHERE oauth_id = $1 LIMIT 1`, guser.UserID)
	if errors.Is(err, sql.ErrNoRows) {

        // Create new user
        _, err = u.db.Exec(`INSERT INTO users (oauth_id, email, first_name, last_name, date_joined) VALUES ($1, $2, $3, $4, $5)`, guser.UserID, guser.Email, guser.FirstName, guser.LastName, time.Now())
        if err != nil {
            return nil, err
        }

        // Im WORKING HERE NOW


	}
	return &user, err
}
