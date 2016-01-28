package main

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/informatik-mannheim/iws_golang/iwsimage"
	"github.com/informatik-mannheim/iws_golang/viewModel"
	"github.com/starmanmartin/simple-router"
	"github.com/starmanmartin/simple-router/view"
)

var (
	indexTemplate, notFoundTemplate, showTemplate *template.Template
	publicPath                                    string
)

func getIndex(w http.ResponseWriter, r *router.Request) (bool, error) {
	indexTemplate.ExecuteTemplate(w, "base", viewModel.NewIndex("Starman"))
	return false, nil
}

func onNotFound(w http.ResponseWriter, r *router.Request) {
	notFoundTemplate.ExecuteTemplate(w, "base", viewModel.NewNotFound("Not found"))
}

func onError(err error, w http.ResponseWriter, r *router.Request) {
	indexTemplate.ExecuteTemplate(w, "base", viewModel.NewIndex(err.Error()))
}

func uploadImage(w http.ResponseWriter, r *router.Request) (bool, error) {
	r.ParseForm()

	filter := r.Form.Get("filter")

	image, has := r.Files["image"]

	if has && strings.HasSuffix(image.Mime, "bmp") {

	}

	return false, errors.New("No image Filter: " + filter)
}

func runWs(port string) {
    // Set view Folder Name
	view.ViewPath = "view"
    
    // Compile html templates
	indexTemplate = view.ParseTemplate("index", "index.html")
	showTemplate = view.ParseTemplate("show", "show.html")
	notFoundTemplate = view.ParseTemplate("404", "404.html")

    // Set error handler
	router.ErrorHandler = onError
    // Set handler for unknown route
	router.NotFoundHandler = onNotFound
	app := router.NewRouter()
    // Add post route and handler
	app.Post("/*", app.UploadPath("upload", false))
    // Add get route and handler
	app.Get("/", getIndex)
    
	publicPath = app.Public("/public")

	http.ListenAndServe(port, app)
}
