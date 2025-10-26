package domain

import "time"

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Meter struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Label     string    `db:"label" json:"label"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Reading struct {
	ID         string    `db:"id" json:"id"`
	MeterID    string    `db:"meter_id" json:"meter_id"`
	ReadingKwh float64   `db:"reading_kwh" json:"reading_kwh"`
	ReadingAt  time.Time `db:"reading_at" json:"reading_at"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// MeterPhoto represents a photo (faulty/fixed) associated with a meter.
type MeterPhoto struct {
	ID        string    `db:"id" json:"id"`
	MeterID   string    `db:"meter_id" json:"meter_id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Type      string    `db:"type" json:"type"` // "faulty" or "fixed"
	Path      string    `db:"path" json:"path"` // relative file path
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
