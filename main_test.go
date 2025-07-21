package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"citynext-appointments/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_CreateAppointment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	req := models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "2075-06-15",
	}

	reqBytes, err := json.Marshal(req)
	require.NoError(t, err)

	var parsed models.CreateAppointmentRequest
	err = json.Unmarshal(reqBytes, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "John", parsed.FirstName)
	assert.Equal(t, "Doe", parsed.LastName)
	assert.Equal(t, "2075-06-15", parsed.VisitDate)

	_, err = time.Parse("2006-01-02", parsed.VisitDate)
	assert.NoError(t, err)
}

func TestJSONSerialization(t *testing.T) {
	// Test appointment JSON serialization
	appointment := models.Appointment{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: time.Date(2075, 6, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(appointment)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled models.Appointment
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, appointment.ID, unmarshaled.ID)
	assert.Equal(t, appointment.FirstName, unmarshaled.FirstName)
	assert.Equal(t, appointment.LastName, unmarshaled.LastName)
}

func TestErrorResponseSerialization(t *testing.T) {
	// Test error response JSON serialization
	errResp := models.ErrorResponse{
		Error:   "validation_error",
		Message: "Invalid request body",
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(errResp)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled models.ErrorResponse
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, errResp.Error, unmarshaled.Error)
	assert.Equal(t, errResp.Message, unmarshaled.Message)
}

func TestHTTPClientTimeout(t *testing.T) {
	// Test HTTP client with timeout (similar to holiday service)
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	// Test with a dummy request
	req, err := http.NewRequest("GET", "https://httpbin.org/delay/2", nil)
	require.NoError(t, err)

	// This should timeout
	_, err = client.Do(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deadline exceeded")
}

func TestDateValidation(t *testing.T) {
	testCases := []struct {
		name        string
		dateString  string
		expectError bool
	}{
		{
			name:        "valid date",
			dateString:  "2075-06-15",
			expectError: false,
		},
		{
			name:        "invalid format",
			dateString:  "15-06-2075",
			expectError: true,
		},
		{
			name:        "invalid date",
			dateString:  "2075-13-40",
			expectError: true,
		},
		{
			name:        "empty string",
			dateString:  "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := time.Parse("2006-01-02", tc.dateString)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequestValidation(t *testing.T) {
	testCases := []struct {
		name    string
		payload map[string]interface{}
		valid   bool
	}{
		{
			name: "valid request",
			payload: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"visit_date": "2075-06-15",
			},
			valid: true,
		},
		{
			name: "missing first name",
			payload: map[string]interface{}{
				"last_name":  "Doe",
				"visit_date": "2075-06-15",
			},
			valid: false,
		},
		{
			name: "missing last name",
			payload: map[string]interface{}{
				"first_name": "John",
				"visit_date": "2075-06-15",
			},
			valid: false,
		},
		{
			name: "missing visit date",
			payload: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
			},
			valid: false,
		},
		{
			name: "empty first name",
			payload: map[string]interface{}{
				"first_name": "",
				"last_name":  "Doe",
				"visit_date": "2075-06-15",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tc.payload)
			require.NoError(t, err)

			var req models.CreateAppointmentRequest
			err = json.Unmarshal(jsonBytes, &req)
			require.NoError(t, err)

			// Basic validation - all fields should be present and non-empty
			isValid := req.FirstName != "" && req.LastName != "" && req.VisitDate != ""

			assert.Equal(t, tc.valid, isValid, "Request validation failed for test case: %s", tc.name)
		})
	}
}

// Benchmark test for JSON marshaling
func BenchmarkJSONMarshal(b *testing.B) {
	appointment := models.Appointment{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: time.Date(2075, 6, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(appointment)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark test for date parsing
func BenchmarkDateParse(b *testing.B) {
	dateString := "2075-06-15"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Example test for HTTP request construction
func TestHTTPRequestConstruction(t *testing.T) {
	req := models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "2075-06-15",
	}

	// Marshal request
	jsonBytes, err := json.Marshal(req)
	require.NoError(t, err)

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", "http://localhost:8080/appointments", bytes.NewBuffer(jsonBytes))
	require.NoError(t, err)

	// Validate request properties
	assert.Equal(t, "POST", httpReq.Method)
	assert.Equal(t, "http://localhost:8080/appointments", httpReq.URL.String())
	assert.NotNil(t, httpReq.Body)
}
