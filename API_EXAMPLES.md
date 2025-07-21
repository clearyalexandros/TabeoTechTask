# CityNext Appointments API Examples

This file contains example API requests for testing the CityNext Appointments API.

## Valid Appointment Request

```bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "visit_date": "2075-06-15"
  }'
```

Expected Response (201 Created):
```json
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "visit_date": "2075-06-15T00:00:00Z",
  "created_at": "2075-01-01T10:00:00Z"
}
```

## Try to Book on UK Public Holiday (via Nager.Date API)

```bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "visit_date": "2075-12-25"
  }'
```

Expected Response (400 Bad Request):
```json
{
  "error": "public_holiday",
  "message": "Cannot book appointment on a public holiday"
}
```

## Try to Book Duplicate Appointment

```bash
# First request (should succeed)
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Alice",
    "last_name": "Johnson",
    "visit_date": "2075-07-01"
  }'

# Second request for same date (should fail)
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Bob",
    "last_name": "Wilson",
    "visit_date": "2075-07-01"
  }'
```

Expected Response for duplicate (409 Conflict):
```json
{
  "error": "duplicate_appointment",
  "message": "appointment already exists for date 2075-07-01"
}
```

## Try to Book in the Past

```bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Charlie",
    "last_name": "Brown",
    "visit_date": "2020-01-01"
  }'
```

Expected Response (400 Bad Request):
```json
{
  "error": "past_date",
  "message": "visit date cannot be in the past"
}
```

## Invalid Request Format

```bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "David",
    "visit_date": "2075-08-15"
  }'
```

Expected Response (400 Bad Request):
```json
{
  "error": "validation_error",
  "message": "Invalid request body: Key: 'CreateAppointmentRequest.LastName' Error:Tag: 'required'"
}
```

## Invalid Date Format

```bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Eve",
    "last_name": "Davis",
    "visit_date": "15-06-2075"
  }'
```

Expected Response (400 Bad Request):
```json
{
  "error": "invalid_date",
  "message": "Invalid date format, expected YYYY-MM-DD"
}
```

## PowerShell Examples

For Windows PowerShell users:

```powershell
# Valid appointment
Invoke-RestMethod -Uri "http://localhost:8080/appointments" -Method Post -ContentType "application/json" -Body '{
    "first_name": "John",
    "last_name": "Doe",
    "visit_date": "2075-06-15"
}'

# Public holiday (should fail)
Invoke-RestMethod -Uri "http://localhost:8080/appointments" -Method Post -ContentType "application/json" -Body '{
    "first_name": "Jane",
    "last_name": "Smith",
    "visit_date": "2075-12-25"
}'
```
