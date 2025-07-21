package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"citynext-appointments/internal/constants"
	"citynext-appointments/internal/models"
)

type HolidayService struct {
	client *http.Client
}

func NewHolidayService() *HolidayService {
	return &HolidayService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *HolidayService) IsPublicHoliday(ctx context.Context, date time.Time) (bool, error) {
	year := date.Year()
	url := fmt.Sprintf(constants.NagerDateAPIURL, year)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to fetch public holidays: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var holidays []models.PublicHoliday
	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	dateStr := date.Format(constants.DateLayout)
	for _, holiday := range holidays {
		if holiday.Date == dateStr {
			return true, nil
		}
	}

	return false, nil
}
