package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Axontik/comin-time-service/internal/handler"
	"github.com/Axontik/comin-time-service/internal/middleware"
	"github.com/Axontik/comin-time-service/internal/repository"
	"github.com/Axontik/comin-time-service/internal/service"
	"github.com/Axontik/comin-time-service/pkg/auth"

	// "github.com/Axontik/comin-time-service/pkg/employee"
	"github.com/Axontik/comin-time-service/pkg/organization"
)

type Application struct {
	db          *gorm.DB
	timeHandler *handler.TimeHandler
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	app := &Application{}

	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	app.db = db

	// Initialize dependencies
	app.initializeDependencies()

	// Setup router
	router := setupRouter(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Using 8084 for time service
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://comin_owner:Ye5rfjcIB7FX@ep-flat-shadow-a8onelva.eastus2.azure.neon.tech/comin?sslmode=require"
	}

	// Run migrations
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize migrations: %v", err)
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Warning: Failed to run migrations: %v", err)
		}
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	return gorm.Open(postgres.Open(dbURL), config)
}

func (app *Application) initializeDependencies() {
	// Initialize repositories
	timeRepo := repository.NewTimeRepository(app.db)

	// Initialize services
	timeService := service.NewTimeService(timeRepo)

	// Initialize handlers
	app.timeHandler = handler.NewTimeHandler(timeService)
}

func (app *Application) healthHandler(c *gin.Context) {
	// Check DB connection
	sqlDB, err := app.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "reason": "database connection error"})
		return
	}
	err = sqlDB.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy", "reason": "database ping failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

func setupRouter(app *Application) *gin.Engine {
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8080/api/v1/auth"
	}
	authClient := auth.NewAuthClient(authServiceURL)

	orgServiceURL := os.Getenv("ORGANIZATION_SERVICE_URL")
	if orgServiceURL == "" {
		orgServiceURL = "http://localhost:8081/api/v1"
	}
	orgClient := organization.NewOrganizationClient(orgServiceURL)

	// employeeClient := employee.NewEmployeeClient(os.Getenv("EMPLOYEE_SERVICE_URL"))
	// if employeeClient == nil {
	// 	employeeClient = employee.NewEmployeeClient("http://localhost:8082/api/v1")
	// }

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	// Health check
	router.GET("/health", app.healthHandler)

	// API routes
	api := router.Group("/api/v1/time")
	{
		// QR Code routes
		qrCodes := api.Group("/organizations/:organization_id/qr-codes")
		qrCodes.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			qrCodes.POST("/", app.timeHandler.GenerateQRCode)
			qrCodes.GET("/:employee_id", app.timeHandler.GetEmployeeQRCodes)
		}

		// Attendance routes (public for check-in/check-out)
		attendance := api.Group("/attendance")
		{
			attendance.POST("/check-in", app.timeHandler.CheckIn)
			attendance.POST("/check-out", app.timeHandler.CheckOut)
		}

		// Protected attendance routes
		orgAttendance := api.Group("/organizations/:organization_id/attendance")
		orgAttendance.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			orgAttendance.GET("/", app.timeHandler.ListAttendances)
			// orgAttendance.GET("/:employee_id", app.timeHandler.GetEmployeeAttendance)
			// orgAttendance.GET("/:employee_id/summary", app.timeHandler.GetAttendanceSummary)
		}

		// Timesheet routes
		timesheets := api.Group("/organizations/:organization_id/employees/:employee_id/timesheets")
		timesheets.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			timesheets.POST("/", app.timeHandler.CreateTimesheet)
			// timesheets.GET("/", app.timeHandler.ListTimesheets)
			timesheets.GET("/:id", app.timeHandler.GetTimesheet)
			timesheets.PUT("/:id", app.timeHandler.UpdateTimesheet)
			timesheets.DELETE("/:id", app.timeHandler.DeleteTimesheet)
			// timesheets.PUT("/:id/approve", app.timeHandler.ApproveTimesheet)
			// timesheets.PUT("/:id/reject", app.timeHandler.RejectTimesheet)
		}

		// Reports
		reports := api.Group("/organizations/:organization_id/reports")
		reports.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			// reports.GET("/attendance", app.timeHandler.GetAttendanceReport)
			// reports.GET("/timesheets", app.timeHandler.GetTimesheetReport)
		}
	}

	return router
}
