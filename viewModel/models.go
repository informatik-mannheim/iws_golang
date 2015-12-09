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

//NewIndex returns a new Instance of index
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

//NewShow returns a new Instance of show
func NewShow(sampleString, href, filter string) (show *Show){
	show = &Show{}
	show.Main.SampleString = sampleString
	show.Main.Href = href
	show.Main.Filter = filter
	return 
}