package branch

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/dechristopher/dchr.host/src/common"
	"github.com/gofiber/fiber/v2"
)

// Handler executes the branch calculator page template
func Handler(c *fiber.Ctx) error {
	return common.HandleTemplate(c, "branch",
		"Branch Calculator", nil, 200)
}

// CalcHandler executes the branch calculator page template
func CalcHandler(c *fiber.Ctx) error {
	// ensure valid submission
	if c.FormValue("calc") == "" {
		return c.Redirect("/branch#oops", http.StatusFound)
	}

	o := c.FormValue("origin")
	if o == "" {
		return c.Redirect("/branch#error", http.StatusFound)
	}

	bid, err := strconv.Atoi(o)
	if err != nil {
		return c.Redirect("/branch#error", http.StatusFound)
	}

	origin := Get(bid)
	if origin == NoBranch {
		return c.Redirect("/branch#error", http.StatusFound)
	}

	var branches Branches

	// assemble destination branches
	destNum := 1
	for {
		key := fmt.Sprintf("d%d", destNum)
		destSelection := c.FormValue(key)
		if destSelection == "" {
			break
		}

		dbId, err := strconv.Atoi(destSelection)

		if err != nil {
			return c.Redirect("/branch#oops", http.StatusFound)
		}

		dest := Get(dbId)
		if dest == NoBranch {
			return c.Redirect("/branch#error", http.StatusFound)
		}

		branches = append(branches, dest)

		destNum++
	}

	var calc Calculation

	if len(branches) != 0 {
		log.Printf("%+v", branches)

		// calculate result
		dist := origin.RoundTrip(branches...)

		calc = Calculation{
			Distance:     dist,
			Cost:         math.Ceil(((float64(dist) * FedMileage2021) * 100) / 100),
			Time:         math.Ceil(((float64(dist) / float64(AverageSpeedLimit)) * 100) / 100),
			Origin:       origin,
			Destinations: branches,
		}

		log.Printf("%+v", calc)
	}

	return common.HandleTemplate(c, "branch", "Branch Calculator", calc, 200)
}
