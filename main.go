package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/dechristopher/dchr.host/src/common"
	"github.com/dechristopher/dchr.host/src/www"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

var (
	port string

	//go:embed static/template/*
	templates embed.FS
	//go:embed static/res/*
	resources embed.FS

	templateFs http.FileSystem
	resourceFs http.FileSystem

	// fiber html template engine
	engine *html.Engine
)

// main does the thing
func main() {
	_ = godotenv.Load()
	port = os.Getenv("PORT")

	log.Printf("DCHR.HOST - :%s - %s", port, common.GetEnv())

	// make filesystem location decision based on environment
	templateFs = common.PickFS(!common.IsProd(), templates, "./static/template")
	resourceFs = common.PickFS(!common.IsProd(), resources, "./static/res")
	// populate template engine from templates filesystem
	engine = html.NewFileSystem(templateFs, ".html")

	// enable template engine reloading on dev
	engine.Reload(!common.IsProd())

	// add custom incrementer template function
	engine.AddFunc("inc", func(i int) int {
		return i + 1
	})

	r := fiber.New(fiber.Config{
		ServerHeader:          "dchr.host",
		CaseSensitive:         true,
		ErrorHandler:          nil,
		DisableStartupMessage: true,
		Views:                 engine,
	})

	www.WireHandlers(r, resourceFs)

	// Graceful shutdown with SIGINT
	// SIGTERM and others will hard kill
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		_ = r.Shutdown()
	}()

	// listen for connections on primary listening port
	if err := r.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Println(err)
	}

	// Exit cleanly
	log.Printf("DCHR.HOST - shutdown")
	os.Exit(0)
}
