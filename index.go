package main

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/informatik-mannheim/iws_golang/iwsimage"
	"github.com/informatik-mannheim/iws_golang/viewModel"
	"github.com/starmanmartin/simple-router"
	"github.com/starmanmartin/simple-router/view"
)

var (
	indexTemplate, showTemplate *template.Template
	publicPath                  string
)

func getIndex(w http.ResponseWriter, r *router.Request) (bool, error) {
	indexTemplate.ExecuteTemplate(w, "base", viewModel.NewIndex("Starman"))
	return false, nil
}

func onError(err error, w http.ResponseWriter, r *router.Request) {
	indexTemplate.ExecuteTemplate(w, "base", viewModel.NewIndex(err.Error()))
}

func uploadImage(w http.ResponseWriter, r *router.Request) (bool, error) {
	r.ParseForm()

	filter := r.Form.Get("filter")

	image, has := r.Files["image"]

	if has && strings.HasSuffix(image.Mime, "bmp") {
		iData := iwsimage.NewImageData()
		iData.LoadFile(filepath.Join(image.Path, image.Name))
		switch filter {
		case "gray":
			iData.Filter(iwsimage.GrayFilter)
		default:
			iData.Filter(iwsimage.BlueFilterGenerator(10))
		}

		iData.SaveFile(filepath.Join(publicPath, "img", image.Name))

		showTemplate.ExecuteTemplate(w, "base", viewModel.NewShow("Martin", "public/img/"+image.Name, filter))
		return false, nil
	}

	indexTemplate.ExecuteTemplate(w, "base", viewModel.NewIndex("Starman"))
	return false, errors.New("No image")
}

func main() {
	view.ViewPath = "view"
	indexTemplate = view.ParseTemplate("index", "index.html")
	showTemplate = view.ParseTemplate("show", "show.html")

	router.ErrorHandler = onError
	app := router.NewRouter()
	app.Post("/*", app.UploadPath("upload", false))
	app.Get("/", getIndex)
	app.Post("/upload", uploadImage)
	publicPath = app.Public("/public")

	http.ListenAndServe(":8080", app)
}
