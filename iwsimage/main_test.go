package iwsimage

import (
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"testing"
	"strconv"
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

	imgData := NewImageData()
	colorCollection := [...]float64{1.2, 0.7, 0.7, 0.7, 1.8, 0.7, 0.7, 0.7, 1.8, 0.3, 1.6, 1.6}
	if err := imgData.LoadFile(src); err != nil {
		panic(err.Error())
	}
	
	imageChan := make(chan *ImageData, 4)
	for i := 0; i < 4; i++ {
		tmpDest := filepath.Join(filepath.Dir(dest), "img_" + strconv.Itoa(i) + ".bmp")
		go func(filter []float64, isCopy bool) {
			var newImgData *ImageData
			if(isCopy) {
				newImgData = imgData.Copy()
			} else {
				newImgData = imgData;
			}
			newImgData.Filter(GreenFilterGenerator(filter[0]))
			newImgData.Filter(RedFilterGenerator(filter[1]))
			newImgData.Filter(BlueFilterGenerator(filter[2]))
			
			newImgData.SaveFile(tmpDest)
			imageChan <- newImgData
		}(colorCollection[i*3:i*3+3], i<3)
	}
	
	imgData1 := <- imageChan
	
	imgData1.AssembleLeft(<- imageChan)
	imgData2 := <- imageChan
	imgData2.AssembleLeft(<- imageChan)

	imgData1.AssembleTop(imgData2)
	
	imgData1.AssembleOverlayer(imgData1)
	

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
