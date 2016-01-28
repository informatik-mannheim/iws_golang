package iwsimage

import (
	"errors"
	"math"
	_"sync"
)

// FilterFunction defines a function signature for image filters
// gets a concrete image packed in an ImageData struct
// returns error if image is empty or not valid
type FilterFunction func(imgData *ImageData) error


// GrayFilter calculate on each pixel a average
func GrayFilter(imgData *ImageData) error {
	if imgData.Width <= 0 || len(imgData.PixelData) <= 0 {
		return errors.New("No Image Data")
	}

	for i := 0; i < len(imgData.PixelData); i += 3 {
		average := int((imgData.PixelData[i] + imgData.PixelData[i+1] + imgData.PixelData[i+2]) / 3)
		imgData.PixelData[i] = average
		imgData.PixelData[i+1] = average
		imgData.PixelData[i+2] = average
	}

	return nil
}


// OldGreenFilterGenerator returns a filter function for the red Color
func OldGreenFilterGenerator(value float64) FilterFunction {
	return oldColorFilterGenerator(1, value)
}

// OldBlueFilterGenerator returns a filter function for the red blue
func OldBlueFilterGenerator(value float64) FilterFunction {
	return oldColorFilterGenerator(0, value)
}

// OldRedFilterGenerator returns a filterfunction for the red Color
func OldRedFilterGenerator(value float64) FilterFunction {
	return oldColorFilterGenerator(2, value)
}

// Advanced topic:
// you might be able to optimize this function with concurrency
func oldColorFilterGenerator(index int, value float64) FilterFunction {
	return func(imgData *ImageData) error {
		if imgData.Width <= 0 || len(imgData.PixelData) <= 0 {
			return errors.New("No Image Data")
		}

		maxColVal := 255.0
		for i := 0; i < len(imgData.PixelData); i += imgData.Width * 3 {
			start := i+index;
			for l := start; l < (start + 3*imgData.Width); l += 3 {
				colVal := float64(imgData.PixelData[l]) * value
				imgData.PixelData[l] = int(math.Floor(math.Min(maxColVal, colVal)))
			}
		}

		return nil
	}
}

