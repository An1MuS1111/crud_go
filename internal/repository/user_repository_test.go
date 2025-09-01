package repository

import (
	"context"
	"crud/internal/models"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	_ "github.com/lib/pq"
)

var testDB *bun.DB

// func TestMain(m *testing.M) {
// 	ctx := context.Background()

// 	// Start Postgres container
// 	req := testcontainers.ContainerRequest{
// 		Image:        "postgres:15",
// 		ExposedPorts: []string{"5433/tcp"},
// 		Env: map[string]string{
// 			"POSTGRES_PASSWORD": "secret",
// 			"POSTGRES_USER":     "postgres",
// 			"POSTGRES_DB":       "crud_test",
// 		},
// 		WaitingFor: wait.ForListeningPort("5433/tcp"),
// 	}
// 	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
// 		ContainerRequest: req,
// 		Started:          true,
// 	})
// 	if err != nil {
// 		log.Fatalf("could not start postgres container: %v", err)
// 	}

// 	host, _ := pgC.Host(ctx)
// 	port, _ := pgC.MappedPort(ctx, "5433")
// 	dsn := fmt.Sprintf("postgres://testuser:password@%s:%s/testdb?sslmode=disable", host, port.Port())

// 	// Connect using bun
// 	sqldb, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatalf("could not connect to postgres: %v", err)
// 	}

// 	testDB = bun.NewDB(sqldb, pgdialect.New())

// 	// Run schema migration (minimal for testing)
// 	_, err = testDB.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("could not migrate schema: %v", err)
// 	}

// 	code := m.Run()

// 	// Cleanup
// 	_ = pgC.Terminate(ctx)

// 	os.Exit(code)
// }

func TestMain(m *testing.M) {
	ctx := context.Background()

	// DSN for the postgres_test container (Docker Compose)
	// dsn := "postgres://postgres:secret@localhost:5433/crud_test?sslmode=disable"
	dsn := "postgres://postgres:secret@db_test:5432/crud_test?sslmode=disable"

	// Wait for Postgres to be ready
	var sqldb *sql.DB
	var err error
	maxRetries := 10
	for i := range maxRetries {
		sqldb, err = sql.Open("postgres", dsn)
		if err == nil {
			err = sqldb.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Waiting for Postgres to be ready (%d/%d)...", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("could not connect to postgres_test after %d retries: %v", maxRetries, err)
	}

	testDB = bun.NewDB(sqldb, pgdialect.New())

	// Clean up old schema before migration
	_, _ = testDB.NewDropTable().Model((*models.User)(nil)).IfExists().Cascade().Exec(ctx)

	// Run schema migration (minimal for testing)
	_, err = testDB.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatalf("could not migrate schema: %v", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup: drop test tables
	_, _ = testDB.NewDropTable().Model((*models.User)(nil)).IfExists().Cascade().Exec(ctx)

	os.Exit(code)
}

func TestUserRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository(testDB)

	// Create
	user := &models.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	// GetByEmail
	fetched, err := repo.GetByEmail(ctx, "john@example.com")
	require.NoError(t, err)
	require.Equal(t, "John Doe", fetched.Name)

	// Update
	fetched.Name = "Jane Doe"
	err = repo.Update(ctx, fetched)
	require.NoError(t, err)

	updated, err := repo.GetByEmail(ctx, "john@example.com")
	require.NoError(t, err)
	require.Equal(t, "Jane Doe", updated.Name)

	// Delete
	err = repo.Delete(ctx, updated)
	require.NoError(t, err)

	_, err = repo.GetByEmail(ctx, "john@example.com")
	require.Error(t, err) // should fail because deleted
}
