package service

import (
	"context"
	"time"

	"citynext-appointments/internal/models"
)

// AppointmentServiceInterface defines the interface for appointment operations
type AppointmentServiceInterface interface {
	CreateAppointment(ctx context.Context, req *models.CreateAppointmentRequest) (*models.Appointment, error)
}

// HolidayServiceInterface defines the interface for holiday operations
type HolidayServiceInterface interface {
	IsPublicHoliday(ctx context.Context, date time.Time) (bool, error)
}
