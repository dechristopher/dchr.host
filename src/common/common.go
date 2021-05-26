package common

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Env string

const (
	prod Env = "prod"
	dev  Env = "dev"
)

var (
	err error

	funcMap = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	t, _ = template.New("").Funcs(funcMap).
		ParseGlob("static/template/*")
)

// HandleTemplate executes the given template
func HandleTemplate(w http.ResponseWriter, file, name string, data interface{}, code int) {
	// Regen templates for development
	if GetEnv() == dev {
		t, err = template.New("").Funcs(funcMap).
			ParseGlob("static/template/*")
	}
	if err != nil {
		log.Printf("Template parse failed error=%s", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(code)
	errX := t.ExecuteTemplate(w, file, GenPageModel(name, data))
	if errX != nil {
		log.Printf("Template execution failed error=%s", errX.Error())
		http.Error(w, errX.Error(), 500)
	}
}

// StrictFs is a Custom strict filesystem implementation to
// prevent directory listings for resources
type StrictFs struct {
	Fs http.FileSystem
}

// Open only allows existing files to be pulled, not directories
func (sfs StrictFs) Open(path string) (http.File, error) {
	f, err := sfs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

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
		return prod
	}
	return dev
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
