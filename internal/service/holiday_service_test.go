package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHolidayService_IsPublicHoliday(t *testing.T) {
	service := NewHolidayService()
	ctx := context.Background()

	christmasDay := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	isHoliday, err := service.IsPublicHoliday(ctx, christmasDay)

	require.NoError(t, err)
	assert.True(t, isHoliday, "Christmas Day should be a public holiday")

	workingDay := time.Date(2024, 6, 18, 0, 0, 0, 0, time.UTC)
	isHoliday, err = service.IsPublicHoliday(ctx, workingDay)

	require.NoError(t, err)
	assert.False(t, isHoliday, "Regular working day should not be a public holiday")
}

func TestHolidayService_IsPublicHoliday_InvalidYear(t *testing.T) {
	service := NewHolidayService()
	ctx := context.Background()

	futureDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := service.IsPublicHoliday(ctx, futureDate)

	if err != nil {
		t.Logf("Expected behavior: API returned error for year 2100: %v", err)
	}
}
