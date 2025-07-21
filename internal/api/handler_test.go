package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"citynext-appointments/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAppointmentService struct {
	mock.Mock
}

func (m *MockAppointmentService) CreateAppointment(ctx context.Context, req *models.CreateAppointmentRequest) (*models.Appointment, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Appointment), args.Error(1)
}

type MockHolidayService struct {
	mock.Mock
}

func (m *MockHolidayService) IsPublicHoliday(ctx context.Context, date time.Time) (bool, error) {
	args := m.Called(ctx, date)
	return args.Bool(0), args.Error(1)
}

func TestHandler_CreateAppointment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAppointmentService := new(MockAppointmentService)
	mockHolidayService := new(MockHolidayService)

	handler := NewHandler(mockAppointmentService, mockHolidayService)

	visitDate := time.Date(2075, 6, 15, 0, 0, 0, 0, time.UTC)
	mockHolidayService.On("IsPublicHoliday", mock.Anything, visitDate).Return(false, nil)

	expectedAppointment := &models.Appointment{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: visitDate,
		CreatedAt: time.Now(),
	}

	req := &models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "2075-06-15",
	}

	mockAppointmentService.On("CreateAppointment", mock.Anything, req).Return(expectedAppointment, nil)

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/appointments", handler.CreateAppointment)
	router.ServeHTTP(w, request)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Appointment
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAppointment.ID, response.ID)
	assert.Equal(t, expectedAppointment.FirstName, response.FirstName)
	assert.Equal(t, expectedAppointment.LastName, response.LastName)

	mockHolidayService.AssertExpectations(t)
	mockAppointmentService.AssertExpectations(t)
}

func TestHandler_CreateAppointment_PublicHoliday(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAppointmentService := new(MockAppointmentService)
	mockHolidayService := new(MockHolidayService)

	handler := NewHandler(mockAppointmentService, mockHolidayService)

	visitDate := time.Date(2075, 12, 25, 0, 0, 0, 0, time.UTC)
	mockHolidayService.On("IsPublicHoliday", mock.Anything, visitDate).Return(true, nil)

	req := &models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "2075-12-25",
	}

	requestBody, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/appointments", handler.CreateAppointment)
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "public_holiday", response.Error)
	assert.Equal(t, "Cannot book appointment on a public holiday", response.Message)

	mockHolidayService.AssertExpectations(t)
}

func TestHandler_CreateAppointment_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAppointmentService := new(MockAppointmentService)
	mockHolidayService := new(MockHolidayService)

	handler := NewHandler(mockAppointmentService, mockHolidayService)

	invalidReq := map[string]interface{}{
		"first_name": "John",
	}

	requestBody, _ := json.Marshal(invalidReq)
	request := httptest.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/appointments", handler.CreateAppointment)
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
}
