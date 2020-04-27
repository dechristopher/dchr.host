package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type env string

const (
	prod env = "prod"
	dev  env = "dev"
)

var (
	port string
	t, _ = template.ParseGlob("static/template/*")
	err  error

	wait time.Duration
)

// main does the thing
func main() {
	_ = godotenv.Load()
	port = os.Getenv("PORT")

	log.Printf("DCHR.HOST - :%s - %s", port, getEnv())

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")

	// Serve static files from /static/res preventing directory listings
	s := http.StripPrefix("/res",
		http.FileServer(strictFs{http.Dir("./static/res")}))
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
func homeHandler(w http.ResponseWriter, r *http.Request) {
	handleTemplate(w, "index.html", "me", nil, 200)
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleTemplate(w, "404.html", "404", nil, 404)
	})
}

// handleTemplate executes the given template
func handleTemplate(w http.ResponseWriter, file, name string, data interface{}, code int) {
	// Regen templates for development
	if getEnv() == dev {
		t, err = template.ParseGlob("static/template/*")
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
