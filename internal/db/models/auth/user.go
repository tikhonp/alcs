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
	ID           int            `db:"id"`
	Email        sql.NullString `db:"email"`
	PasswordHash sql.NullString `db:"password_hash"`
	FirstName    sql.NullString `db:"first_name"`
	LastName     sql.NullString `db:"last_name"`
	PhoneNumber  sql.NullString `db:"phone_number"`
	TelegramId   sql.NullInt64  `db:"telegram_id"`
	OAuthId      sql.NullString `db:"oauth_id"`
	IsActive     bool           `db:"is_active"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
}

func (u *User) FullName() string {
    return fmt.Sprintf("%s %s", u.FirstName.String, u.LastName.String)
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
	return fmt.Sprintf("User{id: %v, email: %v, firstName: %v, lastName: %v}", u.ID, u.Email, firstName, lastName)
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

	// Find or create user by telegram oauth
	FindOrCreateTelegramUser(telegramID string, username string, firstName string, lastName string, photoURL string) (*User, error)

	// GetHostsForOrganization returns all hosts for organization
	GetHostsForOrganization(organizationID int) ([]User, error)
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
	err := u.db.Get(&user, `SELECT * FROM auth_users WHERE id = $1`, id)
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
        JOIN auth_permissions ap on ap.id = auth_user_permissions.permission_id
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
		`SELECT * FROM auth_users WHERE oauth_id = $1`,
		// `UPDATE auth_user SET last_login = now() WHERE oauth_id = $1 RETURNING *`,
		guser.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		// Create new user
		err = u.db.Get(
			&user,
			`INSERT INTO auth_users (oauth_id, email, first_name, last_name, date_joined, last_login) VALUES ($1, $2, $3, $4, now(), now()) RETURNING *`,
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
		`INSERT INTO auth_users (password_hash, email, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`,
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
	err := u.db.Get(&user, "SELECT * FROM auth_users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	if !user.PasswordHash.Valid {
		return nil, errors.New("user does not set password")
	}
	err = compareEncodedHashAndPassword(user.PasswordHash.String, password)
	if err != nil {
		return nil, err
	}
	return &user.ID, nil
}

// Find a user by their Telegram ID
func (u *users) FindByTelegramID(telegramID string) (*User, error) {
	var user User
	query := `SELECT * FROM auth_users WHERE telegram_id = $1`
	err := u.db.Get(&user, query, telegramID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create a new user with Telegram data
func (u *users) CreateTelegramUser(telegramID string, username string, firstName string, lastName string, photoURL string) (*User, error) {
	query := `
		INSERT INTO auth_users (telegram_id, first_name, last_name, is_active, created_at, updated_at)
		VALUES (:telegram_id, :first_name, :last_name, true, NOW(), NOW())
		RETURNING *`

	params := map[string]interface{}{
		"telegram_id": telegramID,
		"first_name":  firstName,
		"last_name":   lastName,
	}

	var user User
	stmt, err := u.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	err = stmt.Get(&user, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Find or create a user with Telegram OAuth data
func (u *users) FindOrCreateTelegramUser(telegramID string, username string, firstName string, lastName string, photoURL string) (*User, error) {
	// Try to find an existing user
	user, err := u.FindByTelegramID(telegramID)
	if err == nil {
		// User exists
		return user, nil
	}

	// If user doesn't exist, create a new one
	return u.CreateTelegramUser(telegramID, username, firstName, lastName, photoURL)
}

func (u *users) GetHostsForOrganization(organizationID int) ([]User, error) {
	var users []User
	query := `
    SELECT * FROM auth_users WHERE id IN (
        SELECT user_id FROM organization_members WHERE organization_id = $1 AND role = 'host'
    )
    `
	err := u.db.Select(&users, query, organizationID)
	return users, err
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
