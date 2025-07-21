package api

import (
	"net/http"
	"time"

	"citynext-appointments/internal/constants"
	"citynext-appointments/internal/models"
	"citynext-appointments/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	appointmentService service.AppointmentServiceInterface
	holidayService     service.HolidayServiceInterface
}

func NewHandler(appointmentService service.AppointmentServiceInterface, holidayService service.HolidayServiceInterface) *Handler {
	return &Handler{
		appointmentService: appointmentService,
		holidayService:     holidayService,
	}
}

func (h *Handler) CreateAppointment(c *gin.Context) {
	var req models.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   constants.ErrorTypeValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	visitDate, err := time.Parse(constants.DateLayout, req.VisitDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   constants.ErrorTypeInvalidDate,
			Message: constants.ErrInvalidDateFormat,
		})
		return
	}

	// Prevent appointments on UK public holidays
	ctx := c.Request.Context()
	isHoliday, err := h.holidayService.IsPublicHoliday(ctx, visitDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   constants.ErrorTypeHolidayCheck,
			Message: "Failed to check public holidays: " + err.Error(),
		})
		return
	}

	if isHoliday {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   constants.ErrorTypePublicHoliday,
			Message: constants.ErrPublicHoliday,
		})
		return
	}

	appointment, err := h.appointmentService.CreateAppointment(ctx, &req)
	if err != nil {
		status := http.StatusInternalServerError
		errorType := constants.ErrorTypeInternal

		// Map errors to appropriate HTTP responses
		switch {
		case err.Error() == constants.ErrPastDate:
			status = http.StatusBadRequest
			errorType = constants.ErrorTypePastDate
		case len(err.Error()) > len(constants.ErrDuplicateAppointment) &&
			err.Error()[:len(constants.ErrDuplicateAppointment)] == constants.ErrDuplicateAppointment:
			status = http.StatusConflict
			errorType = constants.ErrorTypeDuplicateAppt
		}

		c.JSON(status, models.ErrorResponse{
			Error:   errorType,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}
