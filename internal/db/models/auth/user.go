package auth

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	"github.com/tikhonp/alcs/internal/db/utils"
	"github.com/tikhonp/alcs/internal/util/assert"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id          int            `db:"id"`
	Password    sql.NullString `db:"password"`
	LastLogin   sql.NullTime   `db:"last_login"`
	Email       string         `db:"email"`
	PhoneNumber sql.NullString `db:"phone_number"`
	FirstName   sql.NullString `db:"first_name"`
	LastName    sql.NullString `db:"last_name"`
	IsActive    bool           `db:"is_active"`
	DateJoined  time.Time      `db:"date_joined"`
	OAuthId     sql.NullString `db:"oauth_id"`
}

func (u *User) String() string {
	firstName := "nil"
	if u.FirstName.Valid {
		firstName = u.FirstName.String
	}
	lastName := "nil"
	if u.LastName.Valid {
		lastName = u.LastName.String
	}
	return fmt.Sprintf("User{id: %v, email: %v, firstName: %v, lastName: %v}", u.Id, u.Email, firstName, lastName)
}

// User Actions
type Users interface {

	// GetById returns a user by id
	GetById(id int) (*User, error)

	// IsUserHasPermissions checks if user has specified permissions
	IsUserHasPermissions(userId int, permissions ...Permission) (bool, error)

	// FromOAuth finds user by goth user or creates new one
	FromOAuth(guser *goth.User) (*User, error)

	// CreateSuperAdmin creates new user with superadmin priveleages
	CreateSuperAdmin(email, password, firstName, lastName string) error

	// ValidateUserAuth checks if user with given email and password exsits
	// validates passwords and returns user id otherwise error
	ValidateUserAuth(email, password string) (*int, error)
}

type users struct {
	db          *sqlx.DB
	permissions Permissions
}

func NewUsers(db *sqlx.DB) Users {
	return &users{db: db, permissions: NewPermissions(db)}
}

func (u *users) GetById(id int) (*User, error) {
	var user User
	err := u.db.Get(&user, `SELECT * FROM auth_user WHERE id = $1`, id)
	return &user, err
}

func (u *users) IsUserHasPermissions(userId int, permission ...Permission) (bool, error) {
	var permissionCodes []string
	for _, v := range permission {
		permissionCodes = append(permissionCodes, v.code)
	}
	query, args, err := sqlx.In(
		`
        SELECT 1 
        FROM auth_user_permissions
        JOIN auth_permission ap on ap.id = auth_user_permissions.permission_id
        WHERE user_id = ? AND ap.codename IN (?)
        `,
		userId,
		permissionCodes,
	)
	if err != nil {
		return false, err
	}
	query = u.db.Rebind(query)
	rows, err := u.db.Query(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *users) FromOAuth(guser *goth.User) (*User, error) {
	var user User
	err := u.db.Get(
		&user,
		`UPDATE auth_user SET last_login = now() WHERE oauth_id = $1 RETURNING *`,
		guser.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		// Create new user
		err = u.db.Get(
			&user,
			`INSERT INTO auth_user (oauth_id, email, first_name, last_name, date_joined, last_login) VALUES ($1, $2, $3, $4, now(), now()) RETURNING *`,
			guser.UserID,
			guser.Email,
			utils.NewNullString(guser.FirstName),
			utils.NewNullString(guser.LastName),
		)
		return &user, err
	}

	return &user, err
}

func (u *users) CreateSuperAdmin(email, password, firstName, lastName string) error {
	var userId int
	rows, err := u.db.Query(
		`INSERT INTO auth_user (password, email, first_name, last_name, date_joined) VALUES ($1, $2, $3, $4, now()) RETURNING id`,
		getUserPasswordHash(password),
		email,
		utils.NewNullString(firstName),
		utils.NewNullString(lastName),
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&userId)
	} else {
		return errors.New("failed to get user id as returning from insetion")
	}

	return u.permissions.AddPermissionForUser(userId, SuperAdmin)
}

func (u *users) ValidateUserAuth(email, password string) (*int, error) {
	var user User
	err := u.db.Get(&user, "SELECT * FROM auth_user WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	if !user.Password.Valid {
		return nil, errors.New("user does not set password")
	}
	err = compareEncodedHashAndPassword(user.Password.String, password)
	if err != nil {
		return nil, err
	}
	return &user.Id, nil
}

func getUserPasswordHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		assert.NoError(err, "Failed to create password hash")
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func compareEncodedHashAndPassword(hash, password string) error {
	bytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(bytes, []byte(password))
}
