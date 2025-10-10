package domain

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct{ DB *sqlx.DB }

func NewRepository(db *sqlx.DB) *Repository { return &Repository{DB: db} }

// Users
func (r *Repository) CreateUser(ctx context.Context, email, pwdHash string) (User, error) {
	var u User
	err := r.DB.GetContext(ctx, &u, `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at
	`, email, pwdHash)
	return u, err
}
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := r.DB.GetContext(ctx, &u, `SELECT * FROM users WHERE email=$1`, email)
	return u, err
}

// Meters
func (r *Repository) CreateMeter(ctx context.Context, userID, label string) (Meter, error) {
	var m Meter
	err := r.DB.GetContext(ctx, &m, `
		INSERT INTO meters (user_id, label)
		VALUES ($1, $2)
		RETURNING id, user_id, label, created_at
	`, userID, label)
	return m, err
}
func (r *Repository) ListMeters(ctx context.Context, userID string) ([]Meter, error) {
	var ms []Meter
	err := r.DB.SelectContext(ctx, &ms, `SELECT * FROM meters WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	return ms, err
}

// Readings
func (r *Repository) AddReading(ctx context.Context, meterID string, kwh float64, at time.Time) (Reading, error) {
	var rd Reading
	err := r.DB.GetContext(ctx, &rd, `
		INSERT INTO readings (meter_id, reading_kwh, reading_at)
		VALUES ($1, $2, $3)
		RETURNING id, meter_id, reading_kwh, reading_at, created_at
	`, meterID, kwh, at)
	return rd, err
}
func (r *Repository) UsageSum(ctx context.Context, meterID string, from, to time.Time) (float64, error) {
	var total float64
	err := r.DB.GetContext(ctx, &total, `
		SELECT COALESCE(SUM(reading_kwh),0) FROM readings
		WHERE meter_id=$1 AND reading_at >= $2 AND reading_at < $3
	`, meterID, from, to)
	return total, err
}
