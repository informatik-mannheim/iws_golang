package iws_golang

import (
	"errors"
	"math"
)

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

// GreenFilterGenerator returns a filter function for the red Color
func GreenFilterGenerator(value float64) func(imgData *ImageData) error {
	return colorFilterGenerator(1, value)
}

// BlueFilterGenerator returns a filter function for the red blue
func BlueFilterGenerator(value float64) func(imgData *ImageData) error {
	return colorFilterGenerator(0, value)
}

// RedFilterGenerator returns a filterfunction for the red Color
func RedFilterGenerator(value float64) func(imgData *ImageData) error {
	return colorFilterGenerator(2, value)
}

func colorFilterGenerator(index int, value float64) func(imgData *ImageData) error {
	return func(imgData *ImageData) error {
		if imgData.Width <= 0 || len(imgData.PixelData) <= 0 {
			return errors.New("No Image Data")
		}

		maxColVal := 255.0
		for i := 0; i < len(imgData.PixelData); i += imgData.Height * 3 {
			imgData.workerList.Add(1)
			go func(start int, imgData *ImageData) {
				defer imgData.workerList.Done()
				for l := start; l < (start + 3*imgData.Height); l += 3 {
					colVal := float64(imgData.PixelData[l]) * value
					imgData.PixelData[l] = int(math.Floor(math.Min(maxColVal, colVal)))
				}

			}(i+index, imgData)
		}

		imgData.workerList.Wait()
		return nil
	}
}
