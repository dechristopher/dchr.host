package common

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/andersfylling/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/nickname32/discordhook"
)

type Env string

const (
	Prod Env = "prod"
	Dev  Env = "dev"
)

var wa *discordhook.WebhookAPI

func init() {
	webhookId := snowflake.ParseSnowflakeString(os.Getenv("WEBHOOK_ID"))

	api, err := discordhook.NewWebhookAPI(webhookId, os.Getenv("WEBHOOK_TOKEN"),
		false, nil)
	if err != nil {
		panic(err)
	}

	wa = api
}

// StrictFs is a Custom strict filesystem implementation to
// prevent directory listings for resources
type StrictFs struct {
	Fs http.FileSystem
}

// Open only allows existing files to be pulled, not directories
func (sfs StrictFs) Open(path string) (http.File, error) {
	// url decode path to support encoded characters
	path, err := url.QueryUnescape(path)
	if err != nil {
		log.Printf("StrictFS error: %s, %s", path, err.Error())
		return nil, err
	}

	// trim trailing slashes to avoid invalid path errors
	// in fiber's filesystem middleware
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}

	// open file directly if it exists
	f, err := sfs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	// prevent directory listings, only show index file if any
	s, err := f.Stat()
	if err == nil && s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := sfs.Fs.Open(index); err != nil {
			return nil, err
		}
	}
	return f, nil
}

// GetEnv returns the current environment
func GetEnv() Env {
	if os.Getenv("DEPLOY") == "prod" {
		return Prod
	}
	return Dev
}

// IsProd returns true if running in production
func IsProd() bool {
	return GetEnv() == Prod
}

// SendAlert will POST a webhook to discord for alerting purposes
func SendAlert(page, content string) error {
	_, err := wa.Execute(context.Background(), &discordhook.WebhookExecuteParams{
		Content: "received alert",
		TTS:     false,
		Embeds: []*discordhook.Embed{
			{
				Title:       fmt.Sprintf("%s (%s)", page, string(GetEnv())),
				Description: content,
			},
		},
	}, nil, "")

	return err
}

// HandleTemplate will execute the http template engine
// with the given template, name, data, and status
func HandleTemplate(
	c *fiber.Ctx,
	template string,
	name string,
	data interface{},
	status int,
) error {
	return c.Status(status).Render(
		template,
		GenPageModel(name, data),
		"layouts/main")
}

// PageModel contains runtime information that
// can be used during page template rendering
type PageModel struct {
	Env      Env
	PageName string
	Data     interface{}
}

// GenPageModel generates the global page model
func GenPageModel(name string, data interface{}) PageModel {
	return PageModel{
		Env:      GetEnv(),
		PageName: name,
		Data:     data,
	}
}

// MapHash is a struct containing map page route parameters
type MapHash struct {
	Lat     float64
	Lon     float64
	Zoom    float64
	Pitch   float64
	Bearing float64
}

// PickFS returns either an embedded FS or an on-disk FS for the
// given directory path
func PickFS(useDisk bool, e embed.FS, dir string) http.FileSystem {
	if useDisk {
		log.Printf("PickFS - picked disk: %s", dir)
		return http.Dir(dir)
	}

	efs, err := fs.Sub(e, strings.Trim(dir, "./"))
	if err != nil {
		panic(err)
	}

	log.Printf("PickFS - picked embedded: %s", dir)
	return http.FS(efs)
}

// CorsOrigins returns the proper CORS origin configuration
// for the current environment
func CorsOrigins() string {
	return "http://localhost:1337, " +
		"http://localhost:3000, " +
		"https://dchr.host, " +
		"https://*.dchr.host"
}
