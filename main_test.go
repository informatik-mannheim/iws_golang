package iws_golang

import (
	"log"
	"os"
	"os/user"
	"path"
	"testing"
)

func grayFiler(src, dest string) {
	imgData := NewImageData()
	if err := imgData.LoadFile(src); err != nil {
		panic(err.Error())
	}

	imgData.Filter(GrayFilter)

	if err := imgData.SaveFile(dest); err != nil {
		panic(err.Error())
	}

}

func colorCollage(src, dest string) {

	imgData1 := NewImageData()
	colorCollection := [...]float64{1.2, 0.7, 0.7, 0.7, 1.8, 0.7, 0.7, 0.7, 1.8, 0.3, 1.6, 1.6}
	if err := imgData1.LoadFile(src); err != nil {
		panic(err.Error())
	}
	
	imageChan := make(chan *ImageData, 1)
	for i := 1; i < 4; i++ {
		go func(origin *ImageData, filter []float64) {
			imgData := origin.Copy()
			imgData.Filter(GreenFilterGenerator(filter[0]))
			imgData.Filter(RedFilterGenerator(filter[1]))
			imgData.Filter(BlueFilterGenerator(filter[2]))
			imageChan <- imgData
		}(imgData1, colorCollection[i*3:i*3+3])
	}
	
	imgData1.Filter(GreenFilterGenerator(colorCollection[0]))
	imgData1.Filter(RedFilterGenerator(colorCollection[1]))
	imgData1.Filter(BlueFilterGenerator(colorCollection[2]))
	
	
	
	imgData1.AssembleLeft(<- imageChan)
	imgData2 := <- imageChan
	imgData2.AssembleLeft(<- imageChan)

	imgData1.AssembleTop(imgData2)

	if err := imgData1.SaveFile(dest); err != nil {
		panic(err.Error())
	}
}

func BenchmarkGrayFiler(b *testing.B) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		grayFiler(path.Join(dir, "/test/pictures/bridge.bmp"), path.Join(usr.HomeDir, "Desktop/bridgegray.bmp"))
	}
}

func BenchmarkColorCollage(b *testing.B) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		colorCollage(path.Join(dir, "/test/pictures/bridge.bmp"), path.Join(usr.HomeDir, "Desktop/bridgecollor.bmp"))
	}
}

func TestMain(t *testing.T) {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {

		log.Fatal(err)
	}

	colorCollage(path.Join(dir, "/test/pictures/bridge.bmp"), path.Join(usr.HomeDir, "Desktop/bridgecollor.bmp"))

}
