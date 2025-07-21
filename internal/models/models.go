package models

import "time"

// Appointment represents a citizen's appointment booking
type Appointment struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	VisitDate time.Time `json:"visit_date" db:"visit_date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateAppointmentRequest is the payload for booking a new appointment
type CreateAppointmentRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	VisitDate string `json:"visit_date" binding:"required"`
}

// PublicHoliday represents UK public holiday data from the Nager.Date API
type PublicHoliday struct {
	Date        string   `json:"date"`
	LocalName   string   `json:"localName"`
	Name        string   `json:"name"`
	CountryCode string   `json:"countryCode"`
	Fixed       bool     `json:"fixed"`
	Global      bool     `json:"global"`
	Counties    []string `json:"counties"`
	LaunchYear  int      `json:"launchYear"`
	Types       []string `json:"types"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
