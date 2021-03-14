package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/dechristopher/dchr.host/src/branch"
)

type env string

const (
	prod env = "prod"
	dev  env = "dev"
)

var (
	port string
	err  error

	funcMap = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	t, _ = template.New("").Funcs(funcMap).
		ParseGlob("static/template/*")

	wait time.Duration
)

// main does the thing
func main() {
	_ = godotenv.Load()
	port = os.Getenv("PORT")

	log.Printf("DCHR.HOST - :%s - %s", port, getEnv())

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")

	// branch calculator
	r.HandleFunc("/branch", branchHandler).Methods("GET")
	r.HandleFunc("/branch", branchCalcHandler).Methods("POST")

	//predefined route for favicon at root of domain
	r.HandleFunc("/favicon.ico", faviconHandler)

	//predefined route for Glassworks logo
	r.HandleFunc("/gw.png", gwHandler)

	//predefined route for Glassworks jenkins logo
	r.HandleFunc("/jenkins-gw.png", jgwHandler)

	// Serve static files from /static/res preventing directory listings
	sfs := http.FileServer(strictFs{http.Dir("./static/res")})
	s := http.StripPrefix("/res", sfs)
	r.PathPrefix("/res").Handler(s)

	// Custom 404 page
	r.NotFoundHandler = notFoundHandler()

	srv := &http.Server{
		Handler:      r, // router defined above
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	// Serve public HTTP and Websocket endpoints
	//  in a goroutine so they don't block
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// Graceful shutdown with SIGINT. SIGTERM and others will hard kill
	signal.Notify(c, os.Interrupt)

	<-c // Block until we receive SIGINT

	// Create a deadline to wait for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Gracefully shutdown
	_ = srv.Shutdown(ctx)

	// Exit cleanly
	log.Printf("DCHR.HOST - shutdown")
	os.Exit(0)
}

// homeHandler executes the home page template
func homeHandler(w http.ResponseWriter, _ *http.Request) {
	handleTemplate(w, "index.html", "me", nil, 200)
}

// branchHandler executes the branch calculator page template
func branchHandler(w http.ResponseWriter, _ *http.Request) {
	handleTemplate(w, "branch.html", "Branch Calculator", nil, 200)
}

// branchHandler executes the branch calculator page template
func branchCalcHandler(w http.ResponseWriter, r *http.Request) {
	// parse form
	err = r.ParseForm()
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

	origin := branch.Get(bid)
	if origin == branch.NoBranch {
		http.Redirect(w, r, "/branch#error", http.StatusFound)
		return
	}

	var branches branch.Branches

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

		dest := branch.Get(dbid)
		if dest == branch.NoBranch {
			http.Redirect(w, r, "/branch#oops", http.StatusFound)
			return
		}

		branches = append(branches, dest)

		destNum++
	}

	var calc branch.Calculation

	if len(branches) != 0 {
		log.Printf("%+v", branches)

		// calculate result
		dist := origin.RoundTrip(branches...)

		calc = branch.Calculation{
			Distance: dist,
			Cost: math.Ceil(((float64(dist) *
				branch.FedMileage2021) * 100) / 100),
			Time: math.Ceil(((float64(dist) /
				float64(branch.AverageSpeedLimit)) * 100) / 100),
			Origin:       origin,
			Destinations: branches,
		}

		log.Printf("%+v", calc)
	}

	handleTemplate(w, "branch.html", "Branch Calculator", calc, 200)
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleTemplate(w, "404.html", "404", nil, 404)
	})
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/res/ico/favicon.ico")
}

func gwHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/res/gw.png")
}

func jgwHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/res/jenkins-gw.png")
}

// handleTemplate executes the given template
func handleTemplate(w http.ResponseWriter, file, name string, data interface{}, code int) {
	// Regen templates for development
	if getEnv() == dev {
		t, err = template.New("").Funcs(funcMap).
			ParseGlob("static/template/*")
	}
	if err != nil {
		log.Printf("Template parse failed error=%s", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(code)
	errX := t.ExecuteTemplate(w, file, genPageModel(name, data))
	if errX != nil {
		log.Printf("Template execution failed error=%s", errX.Error())
		http.Error(w, errX.Error(), 500)
	}
}

// strictFs is a Custom strict filesystem implementation to
// prevent directory listings for resources
type strictFs struct {
	fs http.FileSystem
}

// Open only allows existing files to be pulled, not directories
func (sfs strictFs) Open(path string) (http.File, error) {
	f, err := sfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err == nil && s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := sfs.fs.Open(index); err != nil {
			return nil, err
		}
	}
	return f, nil
}

// getEnv returns the current environment
func getEnv() env {
	if os.Getenv("DEPLOY") == "prod" {
		return prod
	}
	return dev
}

// pageModel contains runtime information that
// can be used during page template rendering
type pageModel struct {
	Env      env
	PageName string
	Data     interface{}
}

// genPageModel generates the global page model
func genPageModel(name string, data interface{}) pageModel {
	return pageModel{
		Env:      getEnv(),
		PageName: name,
		Data:     data,
	}
}
