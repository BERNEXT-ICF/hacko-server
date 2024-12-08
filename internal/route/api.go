package route

import (
	integration "hacko-app/internal/integration/oauth2google"
	restUser "hacko-app/internal/module/user/handler/rest"
	restClass "hacko-app/internal/module/class/handler/rest"
	restMaterials "hacko-app/internal/module/materials/handler/rest"
	restModules "hacko-app/internal/module/modules/handler/rest"
	restAssignment "hacko-app/internal/module/assignment/handler/rest"
	"hacko-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(app *fiber.App) {
	var (
		googleOauth = integration.NewOauth2googleIntegration()
		api         = app.Group("/users")
	)

	restUser.NewUserHandler(googleOauth).Register(api)
	restClass.NewClassHandler().Register(api)
	restMaterials.NewMaterialsHandler().Register(api)
	restModules.NewModulesHandler().Register(api)
	restAssignment.NewAssignmentHandler().Register(api)

	// fallback route
	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()                       // get the request method
			path   = c.Path()                         // get the request path
			query  = c.Context().QueryArgs().String() // get all query params
			ua     = c.Get("User-Agent")              // get the request user agent
			ip     = c.IP()                           // get the request IP
		)

		log.Info().
			Str("url", c.OriginalURL()).
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("Route not found.")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found"))
	})
}
