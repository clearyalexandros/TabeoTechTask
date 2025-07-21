package constants

const (
	DateLayout = "2006-01-02"
)

const (
	ErrInvalidDateFormat    = "Invalid date format, expected YYYY-MM-DD"
	ErrPastDate             = "Visit date cannot be in the past"
	ErrDuplicateAppointment = "Appointment already exists for date"
	ErrPublicHoliday        = "Cannot book appointment on a public holiday"
)

const (
	ErrorTypeValidation    = "validation_error"
	ErrorTypeInvalidDate   = "invalid_date"
	ErrorTypePastDate      = "past_date"
	ErrorTypeDuplicateAppt = "duplicate_appointment"
	ErrorTypePublicHoliday = "public_holiday"
	ErrorTypeHolidayCheck  = "holiday_check_failed"
	ErrorTypeInternal      = "internal_error"
)

const (
	NagerDateAPIURL = "https://date.nager.at/api/v3/PublicHolidays/%d/GB"
)
