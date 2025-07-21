# CityNext Appointments API

A simple REST API for booking appointments with the CityNext office in 2075. Built with Go and PostgreSQL.

## What it does

- Lets people book appointments through a single endpoint: `POST /appointments`
- Checks UK public holidays automatically using the Nager.Date API to prevent bookings on holidays
- Prevents double-booking on the same date
- Won't let you book appointments in the past
- Stores appointments in PostgreSQL

## How to run it

The easiest way is with Docker:

```bash
docker-compose up
```

The API will be available at http://localhost:8080

## Docker Setup Details

The project uses Docker Compose to orchestrate multiple services:

**Services included:**
- **PostgreSQL database** (port 5432) - Stores appointment data
- **Go application** (port 8080) - The REST API server
- **pgAdmin** (port 5050) - Database management interface (optional)

**What happens when you run `docker-compose up`:**

1. **Database initialization**: PostgreSQL starts and creates the `appointments` table automatically
2. **Application build**: Go app compiles and connects to the database
3. **Network setup**: Services can communicate using service names (e.g., `postgres:5432`)

**Useful Docker commands:**

```bash
# Start in background
docker-compose up -d

# View logs
docker-compose logs app
docker-compose logs postgres

# Stop everything
docker-compose down

# Rebuild after code changes
docker-compose up --build

# Access database directly
docker-compose exec postgres psql -U user -d citynext_appointments
```

**Database access:**
- **Application**: Uses `postgres:5432` (internal Docker network)
- **pgAdmin**: Visit http://localhost:5050 (user: citynext_user, password: citynext_password)


## Making an appointment

Send a POST request to `/appointments`:

```json
{
  "first_name": "John",
  "last_name": "Doe", 
  "visit_date": "2075-06-15"
}
```

If successful, you'll get back:

```json
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "visit_date": "2075-06-15T00:00:00Z",
  "created_at": "2075-01-01T10:00:00Z"
}
```

## Using Postman

For easier testing, import the included Postman collection:

1. Open Postman
2. Click **Import** 
3. Select `CityNext_Appointments_API.postman_collection.json`
4. The collection includes pre-configured requests for:
   - Valid appointment booking
   - Public holiday testing (Christmas Day)
   - Duplicate booking scenarios
   - Past date validation

Make sure your server is running at `http://localhost:8080` before testing.

## Error responses

- **400**: Invalid date, past date, or public holiday
- **409**: Someone already booked that date
- **500**: Something went wrong on our end

## Testing

```bash
# Run all tests
go test ./...

# See test coverage
go test -cover ./...
```

## Concurrency & Performance Testing

The project includes specialized integration tests to validate concurrent request handling and performance under load.

**Running concurrency tests:**

1. Start the server: `go run main.go`
2. In another terminal: `go test -run TestConcurrent`

**Note**: The concurrency tests automatically clean the database before and after running to ensure consistent results.

**What the tests validate:**

- **Concurrent valid requests**: 10 simultaneous appointments with different dates - all should succeed
- **Duplicate prevention**: 5 requests for the same date - only 1 should succeed, 4 should get 409 Conflict
- **Mixed scenarios**: Valid dates, duplicates, and holidays tested together
- **Performance under load**: 50 concurrent requests with response time analysis

**Expected results:**
- Zero data corruption under concurrent access
- Proper duplicate detection (only one appointment per date)
- Response times typically under 50ms per request
- Successful handling of 50+ concurrent requests

These tests demonstrate the application's production readiness and database consistency under concurrent load.

## Project structure

```
├── main.go              # Starts the application
├── internal/
│   ├── api/             # Handles HTTP requests
│   ├── service/         # Business logic
│   ├── db/              # Database connection
│   └── models/          # Data types
└── docker-compose.yml   # Development setup
```
