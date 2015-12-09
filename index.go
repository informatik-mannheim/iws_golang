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

func main() {
	view.ViewPath = "view"
	indexTemplate = view.ParseTemplate("index", "index.html")
	showTemplate = view.ParseTemplate("show", "show.html")
	notFoundTemplate = view.ParseTemplate("404", "404.html")

	router.ErrorHandler = onError
	router.NotFoundHandler = onNotFound
	app := router.NewRouter()
	app.Post("/*", app.UploadPath("upload", false))
	publicPath = app.Public("/public")
	app.Get("/", getIndex)

	http.ListenAndServe(":8080", app)
}
