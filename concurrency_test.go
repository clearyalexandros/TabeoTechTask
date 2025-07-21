// Integration tests for concurrency and performance
// Requires a running server at localhost:8080
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"testing"
	"time"

	"citynext-appointments/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// cleanupDatabase removes test data from the appointments table
func cleanupDatabase(t *testing.T) {
	t.Helper()

	deleteQuery := `
		DELETE FROM appointments 
		WHERE first_name IN (
			'Alice', 'Bob', 'Charlie', 'Diana', 'Eve', 'Frank', 'Grace', 'Henry', 'Ivy', 'Jack',
			'Valid', 'Past', 'Invalid', 'Holiday'
		)
		OR first_name LIKE 'User%'
		OR first_name LIKE 'Load%' 
		OR first_name LIKE 'Bench%';
	`

	cmd := exec.Command("docker-compose", "exec", "-T", "postgres", "psql",
		"-U", "citynext_user", "-d", "citynext_appointments",
		"-c", deleteQuery)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Database cleanup warning: %v, output: %s", err, string(output))
	} else {
		t.Log("Test data cleaned successfully")
	}
}

func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrency test in short mode")
	}

	cleanupDatabase(t)
	defer cleanupDatabase(t)

	baseURL := "http://localhost:8080"

	// Check server availability
	resp, err := http.Get(baseURL + "/appointments")
	if err != nil {
		t.Skipf("Server not available at %s: %v", baseURL, err)
		return
	}
	resp.Body.Close()

	t.Run("ConcurrentValidRequests", func(t *testing.T) {
		requests := []models.CreateAppointmentRequest{
			{"Alice", "Smith", "2075-09-01"},
			{"Bob", "Johnson", "2075-09-02"},
			{"Charlie", "Brown", "2075-09-03"},
			{"Diana", "Wilson", "2075-09-04"},
			{"Eve", "Davis", "2075-09-05"},
			{"Frank", "Miller", "2075-09-06"},
			{"Grace", "Taylor", "2075-09-07"},
			{"Henry", "Anderson", "2075-09-08"},
			{"Ivy", "Thomas", "2075-09-09"},
			{"Jack", "Jackson", "2075-09-10"},
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		results := make([]bool, len(requests))
		responses := make([]*http.Response, len(requests))

		// Send all requests concurrently
		for i, req := range requests {
			wg.Add(1)
			go func(index int, request models.CreateAppointmentRequest) {
				defer wg.Done()

				jsonBody, err := json.Marshal(request)
				require.NoError(t, err)

				resp, err := http.Post(
					baseURL+"/appointments",
					"application/json",
					bytes.NewBuffer(jsonBody),
				)

				mu.Lock()
				if err == nil && resp.StatusCode == http.StatusCreated {
					results[index] = true
					responses[index] = resp
				} else {
					results[index] = false
					if resp != nil {
						resp.Body.Close()
					}
				}
				mu.Unlock()
			}(i, req)
		}

		wg.Wait()

		successCount := 0
		for i, success := range results {
			if success {
				successCount++
				assert.NotNil(t, responses[i], "Response should not be nil for request %d", i)
				if responses[i] != nil {
					responses[i].Body.Close()
				}
			}
		}

		assert.Equal(t, len(requests), successCount, "All concurrent requests should succeed")
	})

	t.Run("ConcurrentDuplicateRequests", func(t *testing.T) {
		duplicateDate := "2075-09-15"
		requestCount := 5

		requests := make([]models.CreateAppointmentRequest, requestCount)
		for i := 0; i < requestCount; i++ {
			requests[i] = models.CreateAppointmentRequest{
				FirstName: fmt.Sprintf("User%d", i),
				LastName:  "Concurrent",
				VisitDate: duplicateDate,
			}
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		successCount := 0
		failureCount := 0

		for i, req := range requests {
			wg.Add(1)
			go func(index int, request models.CreateAppointmentRequest) {
				defer wg.Done()

				jsonBody, err := json.Marshal(request)
				require.NoError(t, err)

				resp, err := http.Post(
					baseURL+"/appointments",
					"application/json",
					bytes.NewBuffer(jsonBody),
				)

				if err == nil {
					defer resp.Body.Close()
					mu.Lock()
					if resp.StatusCode == http.StatusCreated {
						successCount++
					} else {
						failureCount++
					}
					mu.Unlock()
				} else {
					mu.Lock()
					failureCount++
					mu.Unlock()
				}
			}(i, req)
		}

		wg.Wait()

		assert.Equal(t, 1, successCount, "Only one request should succeed for duplicate date")
		assert.Equal(t, requestCount-1, failureCount, "Other requests should fail")
	})

	t.Run("ConcurrentMixedRequests", func(t *testing.T) {
		requests := []struct {
			req            models.CreateAppointmentRequest
			expectedStatus int
		}{
			{models.CreateAppointmentRequest{"Valid", "User1", "2075-09-20"}, http.StatusCreated},
			{models.CreateAppointmentRequest{"Valid", "User2", "2075-09-21"}, http.StatusCreated},
			{models.CreateAppointmentRequest{"Past", "Date", "2075-07-10"}, http.StatusBadRequest},
			{models.CreateAppointmentRequest{"Invalid", "Format", "25-09-2075"}, http.StatusBadRequest},
			{models.CreateAppointmentRequest{"", "Empty", "2075-09-22"}, http.StatusBadRequest},
			{models.CreateAppointmentRequest{"Holiday", "Test", "2075-12-25"}, http.StatusBadRequest},
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		results := make(map[int]int)

		for _, req := range requests {
			results[req.expectedStatus] = 0
		}

		for i, testCase := range requests {
			wg.Add(1)
			go func(index int, tc struct {
				req            models.CreateAppointmentRequest
				expectedStatus int
			}) {
				defer wg.Done()

				jsonBody, err := json.Marshal(tc.req)
				require.NoError(t, err)

				resp, err := http.Post(
					baseURL+"/appointments",
					"application/json",
					bytes.NewBuffer(jsonBody),
				)

				if err == nil {
					defer resp.Body.Close()
					mu.Lock()
					results[resp.StatusCode]++
					mu.Unlock()
				}
			}(i, testCase)
		}

		wg.Wait()

		assert.Equal(t, 2, results[http.StatusCreated], "Should have 2 successful requests")
		assert.True(t, results[http.StatusBadRequest] >= 4, "Should have at least 4 bad requests")
	})
}

func TestServerPerformanceUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	cleanupDatabase(t)
	defer cleanupDatabase(t)

	baseURL := "http://localhost:8080"

	resp, err := http.Get(baseURL + "/appointments")
	if err != nil {
		t.Skipf("Server not available at %s: %v", baseURL, err)
		return
	}
	resp.Body.Close()

	numRequests := 50
	timeout := 30 * time.Second

	t.Run("HighConcurrencyLoad", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		successCount := 0
		errorCount := 0
		startTime := time.Now()

		for i := 0; i < numRequests; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				req := models.CreateAppointmentRequest{
					FirstName: fmt.Sprintf("Load%d", index),
					LastName:  "Test",
					VisitDate: fmt.Sprintf("2075-10-%02d", (index%28)+1),
				}

				jsonBody, err := json.Marshal(req)
				require.NoError(t, err)

				client := &http.Client{Timeout: timeout}
				resp, err := client.Post(
					baseURL+"/appointments",
					"application/json",
					bytes.NewBuffer(jsonBody),
				)

				mu.Lock()
				if err == nil {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusCreated {
						successCount++
					} else {
						errorCount++
					}
				} else {
					errorCount++
				}
				mu.Unlock()
			}(i)
		}

		wg.Wait()
		duration := time.Since(startTime)

		t.Logf("Processed %d requests in %v", numRequests, duration)
		t.Logf("Success: %d, Errors: %d", successCount, errorCount)
		t.Logf("Average response time: %v per request", duration/time.Duration(numRequests))

		assert.True(t, successCount > numRequests/2, "At least half of the requests should succeed")
		assert.True(t, duration < timeout, "All requests should complete within timeout")

		avgResponseTime := duration / time.Duration(numRequests)
		assert.True(t, avgResponseTime < 1*time.Second, "Average response time should be under 1 second")
	})
}

// BenchmarkConcurrentRequests benchmarks the application under concurrent load
func BenchmarkConcurrentRequests(b *testing.B) {
	baseURL := "http://localhost:8080"

	resp, err := http.Get(baseURL + "/appointments")
	if err != nil {
		b.Skipf("Server not available at %s: %v", baseURL, err)
		return
	}
	resp.Body.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := models.CreateAppointmentRequest{
			FirstName: fmt.Sprintf("Bench%d", i),
			LastName:  "Mark",
			VisitDate: fmt.Sprintf("2075-11-%02d", (i%28)+1),
		}

		jsonBody, _ := json.Marshal(req)
		resp, err := http.Post(
			baseURL+"/appointments",
			"application/json",
			bytes.NewBuffer(jsonBody),
		)

		if err == nil {
			resp.Body.Close()
		}
	}
}
