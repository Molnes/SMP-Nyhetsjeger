package router

import (
	"github.com/Molnes/Nyhetsjeger/internal/api/middlewares"
	"github.com/Molnes/Nyhetsjeger/internal/api/web/handlers"
	"github.com/Molnes/Nyhetsjeger/internal/api/web/handlers/api"
	"github.com/Molnes/Nyhetsjeger/internal/data/users/user_roles"
	"github.com/Molnes/Nyhetsjeger/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// Sets up the router for the web server
// Takes care of grouping routes, setting up middleware and registering handlers.
func SetupRouter(e *echo.Echo) {

	databaseConn := database.DB

	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())

	handlers.RegisterQuizHandlers(e)
	dashboardGroup := e.Group("/dashboard")
	handlers.RegisterDashboardHandlers(dashboardGroup)

	// api routes, requiring authentication
	api_group := e.Group("/api/v1")
	authForce := middlewares.NewAuthenticationMiddleware(false)
	api_group.Use(authForce.EncofreAuthentication)

	quiz_api_group := api_group.Group("/quiz")
	api.RegisterQuizApiHandlers(quiz_api_group)

	// admin api routes, requiring admin
	admin_api_group := api_group.Group("/admin")
	enforceAdminMiddleware :=
		middlewares.NewAuthorizationMiddleware(
			databaseConn,
			[]user_roles.Role{
				user_roles.QuizAdmin,
				user_roles.OrganizationAdmin,
			}, false)
	admin_api_group.Use(enforceAdminMiddleware.EnforceRole)
	api.RegisterAdminApiHandlers(admin_api_group)

	e.Static("/static", "assets")

	// websocket for live updates
	e.GET("/ws", handlers.WebsocketHandler)

	// authentication routes
	auth_group := e.Group("/auth")
	handlers.RegisterAuthHandlers(auth_group)
}
