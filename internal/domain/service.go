package domain

import (
	"context"
	"errors"
	"time"

	// other imports
	"golang.org/x/crypto/bcrypt"
)

type Service struct{ Repo *Repository }

func NewService(r *Repository) *Service { return &Service{Repo: r} }

func (s *Service) Signup(ctx context.Context, email, password string) (User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return s.Repo.CreateUser(ctx, email, string(hash))
}
func (s *Service) VerifyUser(ctx context.Context, email, password string) (User, error) {
	u, err := s.Repo.FindUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return User{}, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *Service) CreateMeter(ctx context.Context, userID, label string) (Meter, error) {
	return s.Repo.CreateMeter(ctx, userID, label)
}
func (s *Service) ListMeters(ctx context.Context, userID string) ([]Meter, error) {
	return s.Repo.ListMeters(ctx, userID)
}
func (s *Service) AddReading(ctx context.Context, meterID string, kwh float64, at time.Time) (Reading, error) {
	return s.Repo.AddReading(ctx, meterID, kwh, at)
}
func (s *Service) Usage(ctx context.Context, meterID string, from, to time.Time) (float64, error) {
	return s.Repo.UsageSum(ctx, meterID, from, to)
}

// SaveMeterPhoto stores meter photo metadata.
func (s *Service) SaveMeterPhoto(ctx context.Context, photo *MeterPhoto) error {
	return s.Repo.SaveMeterPhoto(ctx, photo)
}
