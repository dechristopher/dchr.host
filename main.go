package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dechristopher/dchr.host/src/common"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/dechristopher/dchr.host/src/branch"
)

var (
	port string
	wait time.Duration
)

// main does the thing
func main() {
	_ = godotenv.Load()
	port = os.Getenv("PORT")

	log.Printf("DCHR.HOST - :%s - %s", port, common.GetEnv())

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")

	// branch calculator
	r.HandleFunc("/branch", branch.Handler).Methods("GET")
	r.HandleFunc("/branch", branch.CalcHandler).Methods("POST")

	//predefined route for favicon at root of domain
	r.HandleFunc("/favicon.ico", faviconHandler)

	//predefined route for Glassworks logo
	r.HandleFunc("/gw.png", gwHandler)

	//predefined route for Glassworks jenkins logo
	r.HandleFunc("/jenkins-gw.png", jgwHandler)

	// Serve static files from /static/res preventing directory listings
	sfs := http.FileServer(common.StrictFs{Fs: http.Dir("./static/res")})
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
	common.HandleTemplate(w, "index.html", "me", nil, 200)
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		common.HandleTemplate(w, "404.html", "404", nil, 404)
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
