package www

import (
	"net/http"

	"github.com/dechristopher/dchr.host/src/branch"
	"github.com/dechristopher/dchr.host/src/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// WireHandlers builds all of the websocket and http routes
// into the fiber app context
func WireHandlers(r *fiber.App, resourcesFs http.FileSystem) {
	// recover from panics
	r.Use(recover.New())

	// wire up all middleware components
	Wire(r, resourcesFs)

	// index handler
	r.Get("/", indexHandler)

	// branch calculator handlers
	r.Get("/branch", branch.Handler)
	r.Post("/branch", branch.CalcHandler)

	// Custom 404 page
	NotFound(r)
}

// indexHandler executes the home page template
func indexHandler(c *fiber.Ctx) error {
	return common.HandleTemplate(c, "index",
		"me", nil, 200)
}
