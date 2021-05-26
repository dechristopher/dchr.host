package www

import (
	"net/http"
	"os"

	"github.com/dechristopher/dchr.host/src/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const logFormatProd = "${ip} ${header:x-forwarded-for} ${header:x-real-ip} " +
	"[${time}] ${pid} ${locals:requestid} \"${method} ${path} ${protocol}\" " +
	"${status} ${latency} \"${referrer}\" \"${ua}\"\n"

const logFormatDev = "${ip} [${time}] \"${method} ${path} ${protocol}\" " +
	"${status} ${latency}\n"

// Wire attaches all middleware to the given router
func Wire(r fiber.Router, resources http.FileSystem) {
	r.Use(requestid.New())

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: common.CorsOrigins(),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// STDOUT request logger
	r.Use(logger.New(logger.Config{
		// For more options, see the Config section
		TimeZone:   "local",
		TimeFormat: "2006-01-02T15:04:05-0700",
		Format:     logFormat(),
		Output:     os.Stdout,
	}))

	// Predefined route for favicon at root of domain
	r.Use(favicon.New(favicon.Config{
		File:       "ico/favicon.ico",
		FileSystem: resources,
	}))

	// Serve static files from /static/resources preventing directory listings
	r.Use(filesystem.New(filesystem.Config{
		Root:   common.StrictFs{Fs: resources},
		MaxAge: 86400,
	}))
}

// NotFound wires the final 404 handler after all other
// handlers are defined. Acts as the final fallback.
func NotFound(r *fiber.App) {
	r.Use(func(c *fiber.Ctx) error {
		return common.HandleTemplate(c, "404",
			"404", nil, 404)
	})
}

func logFormat() string {
	if common.IsProd() {
		return logFormatProd
	}
	return logFormatDev
}
