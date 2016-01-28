package main

import (
	"errors"
	"github.com/informatik-mannheim/iws_golang/iwsimage"
	"github.com/informatik-mannheim/iws_golang/viewModel"
	"github.com/starmanmartin/simple-router"
	"github.com/starmanmartin/simple-router/view"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
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


// uploadimage handles POST request and filters image
func uploadImage(w http.ResponseWriter, r *router.Request) (bool, error) {
	r.ParseForm()

	filter := r.Form.Get("filter")

	image, has := r.Files["image"]

	if has && strings.HasSuffix(image.Mime, "bmp") {
		var imageToModify *iwsimage.ImageData
		imageToModify = iwsimage.NewImageData()
		imageToModify.LoadFile(filepath.Join(image.Path, image.Name))

		switch filter {
		case "gray":
			imageToModify.Filter(iwsimage.GrayFilter)	
		case "color":
			imageToModify.Filter(iwsimage.OldGreenFilterGenerator(10))
		}

		imageToModify.SaveFile(filepath.Join(publicPath, "img", image.Name))

		// render template
		showTemplate.ExecuteTemplate(w, "base", viewModel.NewShow(image.Name, "public/img/"+image.Name, filter))
		return false, nil
	}

	return false, errors.New("No image Filter: " + filter)
}

// main function that starts the web server
func runWs(port string) {
	
	serverConfiguration()
	
	// creates a router object
	// a router handles REST-Requests and invokes the corresponding go functions
	app := router.NewRouter()
	app.Post("/*", app.UploadPath("upload", false))
	// invoke getIndex after get request
	app.Get("/", getIndex)

	app.Post("/upload", uploadImage)
	
	// start server
	publicPath = app.Public("/public")
	http.ListenAndServe(port, app)
}

func serverConfiguration() {
	// Set view Folder
	view.ViewPath = "view"

	// Parse go template engine specific html templates
	indexTemplate = view.ParseTemplate("index", "index.html")
	showTemplate = view.ParseTemplate("show", "show.html")
	notFoundTemplate = view.ParseTemplate("404", "404.html")

	// Set error handler
	router.ErrorHandler = onError
	// Set handler for unknown route
	router.NotFoundHandler = onNotFound
}
