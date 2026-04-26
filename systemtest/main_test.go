package systemtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internalhttp "github.com/lwlee2608/go-reference/internal/api/http"
	"github.com/lwlee2608/go-reference/internal/db"
	"github.com/lwlee2608/go-reference/internal/db/sqlc"
	"github.com/lwlee2608/go-reference/systemtest/postgres"
	"github.com/lwlee2608/go-reference/systemtest/tests"
)

func TestSystem(t *testing.T) {
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := "postgres"
	dbHost := "localhost"
	schema := "public"

	ctx := context.Background()

	container, err := postgres.StartPostgres(ctx, dbUser, dbPassword, dbName)
	if err != nil {
		t.Skipf("skipping system test because postgres container failed to start: %v", err)
	}
	defer func() {
		assert.NoError(t, postgres.TerminatePostgres(ctx, container))
	}()

	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, port.Num(), dbName)

	require.NoError(t, db.RunMigrations(dbURL, schema))

	pool, err := db.InitDB(ctx, db.Config{URL: dbURL, Schema: schema})
	require.NoError(t, err)
	defer pool.Close()

	services := &internalhttp.Services{Queries: sqlc.New(pool)}

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	internalhttp.SetupRoute(engine, services)

	t.Run("Health", func(t *testing.T) { tests.TestHealthCheck(t, engine) })
	t.Run("Users", func(t *testing.T) { tests.TestUserCRUD(t, engine) })
}
