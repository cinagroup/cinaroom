// Package integration provides helpers for integration tests that require
// a real database connection. All tests using this package should check for
// the TEST_WITH_DB environment variable.
package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB holds a test database connection and provides cleanup utilities.
type TestDB struct {
	DB     *gorm.DB
	dbName string
}

// NewTestDB connects to the test database. Skips if TEST_WITH_DB is not set.
func NewTestDB(t *testing.T) *TestDB {
	t.Helper()
	if os.Getenv("TEST_WITH_DB") == "" {
		t.Skip("Skipping: set TEST_WITH_DB=1 to run integration tests")
	}

	cfg := config.Load()
	testDBName := cfg.Database.DBName + "_test"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		testDBName, cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return &TestDB{DB: db, dbName: testDBName}
}

// Cleanup truncates all test tables. Call in defer.
func (tdb *TestDB) Cleanup(t *testing.T) {
	t.Helper()
	tables := []string{
		"vm_metrics", "vm_logs", "vm_snapshots", "mounts",
		"openclaw_configs", "remote_access", "ip_whitelists", "remote_logs",
		"login_logs", "sessions", "vms", "users", "system_settings",
	}
	for _, tbl := range tables {
		tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tbl))
	}
}

// Close closes the database connection.
func (tdb *TestDB) Close() {
	if tdb.DB != nil {
		sqlDB, _ := tdb.DB.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}
}

// NewAuthenticatedRouter creates a gin.Engine with JWT auth middleware for
// integration testing.
func NewAuthenticatedRouter(jwtCfg *config.JWTConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(&config.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	r.Use(middleware.JWTAuth(jwtCfg))
	return r
}
