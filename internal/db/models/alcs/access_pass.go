package alcs

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type AccessPassStatus string

const (
	AccessPassStatusPending  AccessPassStatus = "pending"
	AccessPassStatusApproved AccessPassStatus = "approved"
	AccessPassStatusRejected AccessPassStatus = "rejected"
)

type AccessPass struct {
	ID             int              `db:"id"`
	UserID         int              `db:"user_id"`
	OrganizationID int              `db:"organization_id"`
	SecurityPostID sql.NullInt64    `db:"security_post_id"`
	QRCode         sql.NullString   `db:"qr_code"`
	VisitorName    sql.NullString   `db:"visitor_name"`
	VehicleNumber  sql.NullString   `db:"vehicle_number"`
	Purpose        sql.NullString   `db:"purpose"`
	ValidFrom      time.Time        `db:"valid_from"`
	ValidUntil     time.Time        `db:"valid_until"`
	ContactDetails sql.NullString   `db:"contact_details"`
	HostID         int              `db:"host_id"`
	Status         AccessPassStatus `db:"status"`
	CreatedAt      time.Time        `db:"created_at"`
	UpdatedAt      time.Time        `db:"updated_at"`
}

type AccessPasses interface {
	CreateAccessPass(pass *AccessPass) error
	UpdateAccessPassStatus(passID int, status AccessPassStatus) error
	GetApproverTelegramID(passID int) (int64, error)
}

type accessPasses struct {
	db *sqlx.DB
}

func NewAccessPasses(db *sqlx.DB) AccessPasses {
	return &accessPasses{db: db}
}

func (ac *accessPasses) CreateAccessPass(pass *AccessPass) error {
	query := `
		INSERT INTO access_passes (visitor_name, vehicle_number, purpose, valid_from, valid_until, contact_details, host_id, status, created_at, updated_at)
		VALUES (:visitor_name, :vehicle_number, :purpose, :valid_from, :valid_until, :contact_details, :host_id, :status, NOW(), NOW())
		RETURNING id`
	_, err := ac.db.NamedExec(query, pass)
	return err
}

func (ac *accessPasses) UpdateAccessPassStatus(passID int, status AccessPassStatus) error {
	query := `UPDATE access_passes SET status = :status, updated_at = NOW() WHERE id = :id`
	params := map[string]interface{}{
		"id":     passID,
		"status": status,
	}
	_, err := ac.db.NamedExec(query, params)
	return err
}

func (ac *accessPasses) GetApproverTelegramID(passID int) (int64, error) {
	query := `SELECT u.telegram_id FROM users u JOIN access_passes ap ON u.id = ap.host_id WHERE ap.id = $1`
	var telegramID int64
	err := ac.db.Get(&telegramID, query, passID)
	return telegramID, err
}
