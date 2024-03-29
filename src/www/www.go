package www

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dechristopher/dchr.host/src/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// WireHandlers builds all the websocket and http routes
// into the fiber app context
func WireHandlers(r *fiber.App, resourcesFs http.FileSystem) {
	// recover from panics
	r.Use(recover.New())

	// wire up all middleware components
	Wire(r, resourcesFs)

	// index handler
	r.Get("/", indexHandler)

	// board page handler
	r.Get("/board", boardHandler)

	// map handler
	r.Get("/map", mapHandler)
	r.Get("/map/:hash", mapHashHandler)

	// Custom 404 page
	NotFound(r)
}

// indexHandler executes the home page template
func indexHandler(c *fiber.Ctx) error {
	return common.HandleTemplate(c, "index",
		"me", nil, 200)
}

// boardHandler executes the home page template
func boardHandler(c *fiber.Ctx) error {
	ip := c.IP()
	go func() {
		err := common.SendAlert("board", fmt.Sprintf("Board page viewed from IP: %s", ip))
		if err != nil {
			fmt.Printf("failed to post discord webhook err=%s\n", err.Error())
		}
	}()
	return common.HandleTemplate(c, "board",
		"kilter board", nil, 200)
}

// mapHandler executes the map page template
func mapHandler(c *fiber.Ctx) error {
	return common.HandleTemplate(c, "map",
		"map", nil, 200)
}

// mapHashHandler executes the map page template with a provided map hash
func mapHashHandler(c *fiber.Ctx) error {
	rawHash := c.Params("hash", "")

	hash := common.MapHash{}

	rawHash = strings.Replace(rawHash, "@", "", -1)

	parts := strings.Split(rawHash, ",")
	if len(parts) < 2 {
		return mapHandler(c)
	}

	for i, part := range parts {
		switch i {
		case 0:
			lat, err := strconv.ParseFloat(part, 64)
			if err != nil {
				fmt.Printf("invalid lat provided: %s, err: %s", part, err.Error())
				return mapHandler(c)
			}
			hash.Lat = lat
		case 1:
			lon, err := strconv.ParseFloat(part, 64)
			if err != nil {
				fmt.Printf("invalid lon provided: %s, err: %s", part, err.Error())
				return mapHandler(c)
			}
			hash.Lon = lon
		case 2:
			part = strings.Replace(part, "z", "", -1)
			zoom, err := strconv.ParseFloat(part, 64)
			if err != nil {
				fmt.Printf("invalid zoom provided: %s, err: %s", part, err.Error())
				continue
			}
			if zoom >= 0 && zoom <= 19 {
				hash.Zoom = zoom
			}
		case 3:
			pitch, err := strconv.ParseFloat(part, 64)
			if err != nil {
				fmt.Printf("invalid pitch provided: %s, err: %s", part, err.Error())
				continue
			}
			if pitch >= 0 && pitch <= 70 {
				hash.Pitch = pitch
			}
		case 4:
			bearing, err := strconv.ParseFloat(part, 64)
			if err != nil {
				fmt.Printf("invalid bearing provided: %s, err: %s", part, err.Error())
				continue
			}
			if bearing >= -180 && bearing <= 180 {
				hash.Bearing = bearing
			}
		}
	}

	return common.HandleTemplate(c, "map",
		"map", hash, 200)
}
