package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDB_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectPing()

	dbInstance := &DB{sqlDB}

	assert.NotNil(t, dbInstance.DB)
	assert.IsType(t, &sql.DB{}, dbInstance.DB)

	err = dbInstance.Ping()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewDB_InvalidConnectionString(t *testing.T) {
	invalidURLs := []string{
		"",
		"invalid://connection/string",
		"postgres://user@:invalid:port/db",
		"postgres://",
	}

	for _, url := range invalidURLs {
		t.Run(fmt.Sprintf("URL_%s", url), func(t *testing.T) {
			db, err := NewDB(url)

			if err == nil {
				assert.Nil(t, db)
			} else {
				assert.Error(t, err)
				assert.Nil(t, db)
			}
		})
	}
}

func TestNewDB_PingFailure(t *testing.T) {
	sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectPing().WillReturnError(sql.ErrConnDone)

	dbInstance := &DB{sqlDB}

	err = dbInstance.Ping()
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_EmbeddedMethods(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	dbInstance := &DB{sqlDB}

	t.Run("Close", func(t *testing.T) {
		assert.NotNil(t, dbInstance.Close)
	})

	t.Run("Query", func(t *testing.T) {
		mock.ExpectQuery("SELECT 1").
			WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(1))

		rows, err := dbInstance.Query("SELECT 1")
		assert.NoError(t, err)
		assert.NotNil(t, rows)
		rows.Close()
	})

	t.Run("Exec", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM test").
			WillReturnResult(sqlmock.NewResult(0, 1))

		result, err := dbInstance.Exec("DELETE FROM test")
		assert.NoError(t, err)
		assert.NotNil(t, result)

		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)
	})

	t.Run("QueryRow", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM test").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		var count int
		err := dbInstance.QueryRow("SELECT COUNT(*) FROM test").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_ContextMethods(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	dbInstance := &DB{sqlDB}
	ctx := context.Background()

	t.Run("QueryContext", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM appointments").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"))

		rows, err := dbInstance.QueryContext(ctx, "SELECT * FROM appointments")
		assert.NoError(t, err)
		assert.NotNil(t, rows)
		rows.Close()
	})

	t.Run("ExecContext", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO appointments").
			WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := dbInstance.ExecContext(ctx, "INSERT INTO appointments (name) VALUES (?)", "test")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("QueryRowContext", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM appointments WHERE name = \\$1").
			WithArgs("test").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

		var id int
		err := dbInstance.QueryRowContext(ctx, "SELECT id FROM appointments WHERE name = $1", "test").Scan(&id)
		assert.NoError(t, err)
		assert.Equal(t, 42, id)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_TransactionMethods(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	dbInstance := &DB{sqlDB}
	ctx := context.Background()

	t.Run("Begin", func(t *testing.T) {
		mock.ExpectBegin()

		tx, err := dbInstance.Begin()
		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})

	t.Run("BeginTx", func(t *testing.T) {
		mock.ExpectBegin()

		tx, err := dbInstance.BeginTx(ctx, nil)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_PreparedStatements(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	dbInstance := &DB{sqlDB}
	ctx := context.Background()

	t.Run("Prepare", func(t *testing.T) {
		mock.ExpectPrepare("SELECT \\* FROM appointments WHERE id = \\$1")

		stmt, err := dbInstance.Prepare("SELECT * FROM appointments WHERE id = $1")
		assert.NoError(t, err)
		assert.NotNil(t, stmt)
	})

	t.Run("PrepareContext", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO appointments")

		stmt, err := dbInstance.PrepareContext(ctx, "INSERT INTO appointments (name) VALUES ($1)")
		assert.NoError(t, err)
		assert.NotNil(t, stmt)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewDB_RealDatabase_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testDatabaseURL := "postgres://test_user:test_password@localhost:5432/test_db?sslmode=disable"

	t.Run("ValidConnection", func(t *testing.T) {
		db, err := NewDB(testDatabaseURL)
		if err != nil {
			t.Skipf("Skipping integration test: no test database available (%v)", err)
			return
		}
		defer db.Close()

		err = db.Ping()
		assert.NoError(t, err)

		var result int
		err = db.QueryRow("SELECT 1").Scan(&result)
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})
}

func BenchmarkDB_QueryRow(b *testing.B) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(b, err)
	defer sqlDB.Close()

	dbInstance := &DB{sqlDB}

	for i := 0; i < b.N; i++ {
		mock.ExpectQuery("SELECT 1").
			WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(1))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result int
		_ = dbInstance.QueryRow("SELECT 1").Scan(&result)
	}
}
