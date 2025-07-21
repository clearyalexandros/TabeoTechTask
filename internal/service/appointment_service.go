package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"citynext-appointments/internal/constants"
	"citynext-appointments/internal/db"
	"citynext-appointments/internal/models"
)

type AppointmentService struct {
	db           *db.DB
	timeProvider func() time.Time
}

func NewAppointmentService(database *db.DB) *AppointmentService {
	return &AppointmentService{
		db:           database,
		timeProvider: time.Now, // Default to real time
	}
}

// NewAppointmentServiceWithTime creates a service with a custom time provider
func NewAppointmentServiceWithTime(database *db.DB, timeProvider func() time.Time) *AppointmentService {
	return &AppointmentService{
		db:           database,
		timeProvider: timeProvider,
	}
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, req *models.CreateAppointmentRequest) (*models.Appointment, error) {
	visitDate, err := time.Parse(constants.DateLayout, req.VisitDate)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.ErrInvalidDateFormat, err)
	}

	// Prevent scheduling appointments in the past
	now := s.timeProvider()
	if visitDate.Before(now.Truncate(24 * time.Hour)) {
		return nil, fmt.Errorf("%s", constants.ErrPastDate)
	}

	// Prevent duplicate appointments on the same date
	exists, err := s.appointmentExistsForDate(ctx, visitDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing appointments: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("%s %s", constants.ErrDuplicateAppointment, req.VisitDate)
	}

	appointment := &models.Appointment{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		VisitDate: visitDate,
	}

	// Parameterized query ($1, $2, $3) prevents SQL injection
	query := `
		INSERT INTO appointments (first_name, last_name, visit_date)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err = s.db.QueryRowContext(ctx, query, appointment.FirstName, appointment.LastName, appointment.VisitDate).
		Scan(&appointment.ID, &appointment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

func (s *AppointmentService) appointmentExistsForDate(ctx context.Context, date time.Time) (bool, error) {
	// Uses indexed visit_date column for fast lookup, parameterized query prevents SQL injection
	query := `SELECT COUNT(*) FROM appointments WHERE visit_date = $1`
	var count int
	err := s.db.QueryRowContext(ctx, query, date).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *AppointmentService) GetAppointmentByDate(ctx context.Context, date time.Time) (*models.Appointment, error) {
	// Uses indexed visit_date column, $1 parameter safely handles user input
	query := `
		SELECT id, first_name, last_name, visit_date, created_at
		FROM appointments
		WHERE visit_date = $1
	`

	appointment := &models.Appointment{}
	err := s.db.QueryRowContext(ctx, query, date).Scan(
		&appointment.ID,
		&appointment.FirstName,
		&appointment.LastName,
		&appointment.VisitDate,
		&appointment.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return appointment, nil
}
