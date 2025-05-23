package routes

import (
	"davet.link/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Use(logger.New())

	app.Use(middlewares.SessionMiddleware())

	registerWebsiteRoutes(app)
	registerAuthRoutes(app)
	registerDashboardRoutes(app)
	registerPanelRoutes(app)
}
