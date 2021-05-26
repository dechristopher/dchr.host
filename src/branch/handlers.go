package branch

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/dechristopher/dchr.host/src/common"
)

// Handler executes the branch calculator page template
func Handler(w http.ResponseWriter, _ *http.Request) {
	common.HandleTemplate(w, "branch.html", "Branch Calculator", nil, 200)
}

// CalcHandler executes the branch calculator page template
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/branch#oops", http.StatusFound)
		return
	}

	// ensure valid submission
	if r.Form.Get("calc") == "" {
		http.Redirect(w, r, "/branch#oops", http.StatusFound)
		return
	}

	o := r.Form.Get("origin")
	if o == "" {
		http.Redirect(w, r, "/branch#error", http.StatusFound)
		return
	}

	bid, err := strconv.Atoi(o)
	if err != nil {
		http.Redirect(w, r, "/branch#error", http.StatusFound)
		return
	}

	origin := Get(bid)
	if origin == NoBranch {
		http.Redirect(w, r, "/branch#error", http.StatusFound)
		return
	}

	var branches Branches

	// assemble destination branches
	destNum := 1
	for {
		key := fmt.Sprintf("d%d", destNum)
		destSelection := r.Form.Get(key)
		if destSelection == "" {
			break
		}

		dbid, err := strconv.Atoi(destSelection)

		if err != nil {
			http.Redirect(w, r, "/branch#oops", http.StatusFound)
			return
		}

		dest := Get(dbid)
		if dest == NoBranch {
			http.Redirect(w, r, "/branch#oops", http.StatusFound)
			return
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

	common.HandleTemplate(w, "branch.html", "Branch Calculator", calc, 200)
}
