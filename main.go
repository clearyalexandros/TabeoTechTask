package main

import (
	"log"
	"os"
	"time"

	"citynext-appointments/internal/api"
	"citynext-appointments/internal/db"
	"citynext-appointments/internal/service"

	"github.com/gin-gonic/gin"
)

// Variables to store the application start time for consistent 2075 simulation
var (
	appStartTime time.Time
	baseDate2075 time.Time
)

// init sets up the time reference when the application starts
func init() {
	appStartTime = time.Now()
	// Set base date to today's date but in 2075
	baseDate2075 = time.Date(2075, appStartTime.Month(), appStartTime.Day(), 12, 0, 0, 0, time.UTC)
}

// timeProvider2075 returns a time as if we're in the year 2075
func timeProvider2075() time.Time {
	// Calculate elapsed time since app started
	elapsed := time.Since(appStartTime)

	// Return the base 2075 date plus the elapsed time
	return baseDate2075.Add(elapsed)
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://citynext_user:citynext_password@localhost:5432/citynext_appointments?sslmode=disable"
	}

	database, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	appointmentService := service.NewAppointmentServiceWithTime(database, timeProvider2075)
	holidayService := service.NewHolidayService()
	handler := api.NewHandler(appointmentService, holidayService)

	router := gin.Default()
	router.POST("/appointments", handler.CreateAppointment)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
