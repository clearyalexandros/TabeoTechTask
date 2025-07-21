package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"citynext-appointments/internal/constants"
	"citynext-appointments/internal/db"
	"citynext-appointments/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppointmentService_CreateAppointment_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockDB := &db.DB{DB: sqlDB}
	service := NewAppointmentService(mockDB)

	req := &models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "2025-08-15",
	}

	expectedID := 1
	expectedCreatedAt := time.Now()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE visit_date = \$1`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery(`INSERT INTO appointments \(first_name, last_name, visit_date\) VALUES \(\$1, \$2, \$3\) RETURNING id, created_at`).
		WithArgs("John", "Doe", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(expectedID, expectedCreatedAt))

	ctx := context.Background()
	result, err := service.CreateAppointment(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedID, result.ID)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, expectedCreatedAt, result.CreatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAppointmentService_CreateAppointment_PastDate(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockDB := &db.DB{DB: sqlDB}
	service := NewAppointmentService(mockDB)

	req := &models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "2025-07-15",
	}

	ctx := context.Background()
	result, err := service.CreateAppointment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), constants.ErrPastDate)
}

func TestAppointmentService_CreateAppointment_InvalidDateFormat(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockDB := &db.DB{DB: sqlDB}
	service := NewAppointmentService(mockDB)

	req := &models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "15-08-2025",
	}

	ctx := context.Background()
	result, err := service.CreateAppointment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), constants.ErrInvalidDateFormat)
}

func TestAppointmentService_CreateAppointment_DuplicateDate(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockDB := &db.DB{DB: sqlDB}
	service := NewAppointmentService(mockDB)

	req := &models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "2025-08-15",
	}

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE visit_date = \$1`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	ctx := context.Background()
	result, err := service.CreateAppointment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), constants.ErrDuplicateAppointment)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAppointmentService_CreateAppointment_DatabaseError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockDB := &db.DB{DB: sqlDB}
	service := NewAppointmentService(mockDB)

	req := &models.CreateAppointmentRequest{
		FirstName: "John",
		LastName:  "Doe",
		VisitDate: "2025-08-15",
	}

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE visit_date = \$1`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	result, err := service.CreateAppointment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check existing appointments")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAppointmentService_GetAppointmentByDate_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockDB := &db.DB{DB: sqlDB}
	service := NewAppointmentService(mockDB)

	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	expectedCreatedAt := time.Now()

	mock.ExpectQuery(`SELECT id, first_name, last_name, visit_date, created_at FROM appointments WHERE visit_date = \$1`).
		WithArgs(testDate).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "visit_date", "created_at"}).
			AddRow(1, "John", "Doe", testDate, expectedCreatedAt))

	ctx := context.Background()
	result, err := service.GetAppointmentByDate(ctx, testDate)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, testDate, result.VisitDate)
	assert.Equal(t, expectedCreatedAt, result.CreatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAppointmentService_GetAppointmentByDate_NotFound(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockDB := &db.DB{DB: sqlDB}
	service := NewAppointmentService(mockDB)

	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`SELECT id, first_name, last_name, visit_date, created_at FROM appointments WHERE visit_date = \$1`).
		WithArgs(testDate).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	result, err := service.GetAppointmentByDate(ctx, testDate)

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAppointmentService_appointmentExistsForDate(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mockDB := &db.DB{DB: sqlDB}
	service := NewAppointmentService(mockDB)

	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE visit_date = \$1`).
		WithArgs(testDate).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	ctx := context.Background()
	exists, err := service.appointmentExistsForDate(ctx, testDate)

	assert.NoError(t, err)
	assert.True(t, exists)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE visit_date = \$1`).
		WithArgs(testDate).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err = service.appointmentExistsForDate(ctx, testDate)

	assert.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
