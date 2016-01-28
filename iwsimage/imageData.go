package iwsimage

import (
	"math"
)

// ImageData is a struct containing all image handle PixelData PixelData
// is an int array holding all image pixel info. The pixel data
// is stored in following order: (i+0):r (i+1):g (i+2):b and  column first
type ImageData struct {
	PixelData  []int
	Width      int
	Height     int
}

// NewImageData returns new instance of ImageData
func NewImageData() (nid *ImageData) {
	nid = &ImageData{}
	
	return
}

// SetSize resets the image size
func (d *ImageData) SetSize(width, height int) {
	d.Width, d.Height = width, height
	d.PixelData = make([]int, 0, height*width*3)
}

// Copy makes a deep copy of a ImageData
func (d *ImageData) Copy() *ImageData {
	newImgData := NewImageData()
	newImgData.Width, newImgData.Height = d.Width, d.Height

	newImgData.PixelData = make([]int, d.Height*d.Width*3)
	for i, e := range d.PixelData {
		newImgData.PixelData[i] = e
	}

	return newImgData
}

// AssembleTop  assembles a image on top of it self
func (d *ImageData) AssembleTop(a *ImageData) {
	var newHeight = d.Height + a.Height
	var newWidth = int(math.Min(float64(d.Width), float64(a.Width)))
	newPixelArray := make([]int, newWidth*newHeight*3)

	for y := 0; y < a.Height*3; y += 3 {
		for x := 0; x < newWidth*3; x++ {
			newPixelArray[y*newWidth+x] = a.PixelData[y*a.Width+x]
		}
	}

	startPoit := a.Height * newWidth * 3

	for y := 0; y < d.Height*3; y += 3 {
		for x := 0; x < newWidth*3; x++ {
			newPixelArray[y*newWidth+startPoit+x] = d.PixelData[y*d.Width+x]
		}
	}

	d.Width, d.Height = newWidth, newHeight
	d.PixelData = newPixelArray
}

// AssembleLeft  assembles a image on the left of it self
func (d *ImageData) AssembleLeft(a *ImageData) {
	var newWidth = d.Width + a.Width
	var newHeight = int(math.Min(float64(d.Height), float64(a.Height)))
	newPixelArray := make([]int, newWidth*newHeight*3)

	for y := 0; y < newHeight*3; y += 3 {
		for x := 0; x < d.Width*3; x++ {
			newPixelArray[y*newWidth+x] = d.PixelData[y*d.Width+x]
		}
		for x := 0; x < a.Width*3; x++ {
			newPixelArray[y*newWidth+x+d.Width*3] = a.PixelData[y*a.Width+x]
		}
	}

	d.Width, d.Height = newWidth, newHeight
	d.PixelData = newPixelArray
}

func (d *ImageData) addPixel(r, g, b int) {
	d.PixelData = append(d.PixelData, r, g, b)
}

// SaveFile saves the image to a given path
func (d *ImageData) SaveFile(dest string) error {
	return savefile(dest, d)
}

// LoadFile reads a file from disk
func (d *ImageData) LoadFile(src string) error {
	return loadFile(src, d)
}

// Filter runs a Filter at the image data. The filter gets past as parameter
func (d *ImageData) Filter(filter FilterFunction) error {
	return filter(d)
}
