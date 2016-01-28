package viewModel

import()

type base struct {}

//Index is Struct to parse index.html
type Index struct {
	Base base
	Main struct {
		SampleString string
	}
}

//NewIndex returns a new Instance of Index
func NewIndex(sampleString string) (index *Index){
	index = &Index{}
	index.Main.SampleString = sampleString
	return
}

//Show is Struct to parse show.html
type Show struct {
	Base base
	Main struct {
		SampleString, Href, Filter string
	}
}

//NewShow returns a new Instance of Show
func NewShow(sampleString, href, filter string) (show *Show){
	show = &Show{}
	show.Main.SampleString = sampleString
	show.Main.Href = href
	show.Main.Filter = filter
	return 
}

//NotFound is Struct to parse 404.html
type NotFound struct {
	Base base
	Main string
}

//NewNotFound returns a new Instance of NotFound
func NewNotFound(message string) (notFound *NotFound){
	notFound = &NotFound{}
	notFound.Main = message
	return 
}