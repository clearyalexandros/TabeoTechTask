package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppointment_JSONSerialization(t *testing.T) {
	createdAt := time.Date(2025, 7, 20, 14, 30, 0, 0, time.UTC)
	visitDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	appointment := Appointment{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: visitDate,
		CreatedAt: createdAt,
	}

	jsonData, err := json.Marshal(appointment)
	require.NoError(t, err)

	expectedJSON := `{"id":1,"first_name":"John","last_name":"Doe","visit_date":"2025-08-15T00:00:00Z","created_at":"2025-07-20T14:30:00Z"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))

	var unmarshaled Appointment
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, appointment.ID, unmarshaled.ID)
	assert.Equal(t, appointment.FirstName, unmarshaled.FirstName)
	assert.Equal(t, appointment.LastName, unmarshaled.LastName)
	assert.Equal(t, appointment.VisitDate, unmarshaled.VisitDate)
	assert.Equal(t, appointment.CreatedAt, unmarshaled.CreatedAt)
}

func TestCreateAppointmentRequest_JSONSerialization(t *testing.T) {
	request := CreateAppointmentRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		VisitDate: "2025-08-20",
	}

	jsonData, err := json.Marshal(request)
	require.NoError(t, err)

	expectedJSON := `{"first_name":"Jane","last_name":"Smith","visit_date":"2025-08-20"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))

	var unmarshaled CreateAppointmentRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, request.FirstName, unmarshaled.FirstName)
	assert.Equal(t, request.LastName, unmarshaled.LastName)
	assert.Equal(t, request.VisitDate, unmarshaled.VisitDate)
}

func TestPublicHoliday_JSONSerialization(t *testing.T) {
	holiday := PublicHoliday{
		Date:        "2025-12-25",
		LocalName:   "Christmas Day",
		Name:        "Christmas Day",
		CountryCode: "GB",
		Fixed:       true,
		Global:      true,
		Counties:    []string{"ENG", "SCT", "WLS", "NIR"},
		LaunchYear:  2000,
	}

	jsonData, err := json.Marshal(holiday)
	require.NoError(t, err)

	var unmarshaled PublicHoliday
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, holiday.Date, unmarshaled.Date)
	assert.Equal(t, holiday.LocalName, unmarshaled.LocalName)
	assert.Equal(t, holiday.Name, unmarshaled.Name)
	assert.Equal(t, holiday.CountryCode, unmarshaled.CountryCode)
	assert.Equal(t, holiday.Fixed, unmarshaled.Fixed)
	assert.Equal(t, holiday.Global, unmarshaled.Global)
	assert.Equal(t, holiday.Counties, unmarshaled.Counties)
	assert.Equal(t, holiday.LaunchYear, unmarshaled.LaunchYear)
}

func TestErrorResponse_JSONSerialization(t *testing.T) {
	errorResp := ErrorResponse{
		Error:   "validation_error",
		Message: "Invalid request body",
	}

	jsonData, err := json.Marshal(errorResp)
	require.NoError(t, err)

	expectedJSON := `{"error":"validation_error","message":"Invalid request body"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))

	var unmarshaled ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, errorResp.Error, unmarshaled.Error)
	assert.Equal(t, errorResp.Message, unmarshaled.Message)
}

func TestAppointment_ZeroValues(t *testing.T) {
	var appointment Appointment

	assert.Equal(t, 0, appointment.ID)
	assert.Equal(t, "", appointment.FirstName)
	assert.Equal(t, "", appointment.LastName)
	assert.True(t, appointment.VisitDate.IsZero())
	assert.True(t, appointment.CreatedAt.IsZero())
}

func TestCreateAppointmentRequest_ZeroValues(t *testing.T) {
	var request CreateAppointmentRequest

	assert.Equal(t, "", request.FirstName)
	assert.Equal(t, "", request.LastName)
	assert.Equal(t, "", request.VisitDate)
}
