package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"
	"strconv"
	"github.com/informatik-mannheim/iws_golang/iwsimage"
)

var numberOfRuns = 3

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	
	timeNow := time.Now()
	for i := 0; i < numberOfRuns; i++ {
		colorCollage(path.Join(dir, "/test/pictures/bridge.bmp"), path.Join(usr.HomeDir, "Desktop/bridgecollor.bmp"))
	}
	
	fmt.Println("average runtime(ns):");
	fmt.Println(time.Since(timeNow).Nanoseconds() / int64(numberOfRuns))
}

func colorCollage(src, dest string) {

	imgData := iwsimage.NewImageData()
	colorCollection := [...]float64{1.2, 0.7, 0.7, 0.7, 1.8, 0.7, 0.7, 0.7, 1.8, 0.3, 1.6, 1.6}
	if err := imgData.LoadFile(src); err != nil {
		panic(err.Error())
	}
	
	imageChan := make(chan *iwsimage.ImageData, 4)
	for i := 0; i < 4; i++ {
		tmpDest := filepath.Join(filepath.Dir(dest), "img_" + strconv.Itoa(i) + ".bmp")
		go func(filter []float64, isCopy bool) {
			var newImgData *iwsimage.ImageData
			if(isCopy) {
				newImgData = imgData.Copy()
			} else {
				newImgData = imgData;
			}
			newImgData.Filter(iwsimage.GreenFilterGenerator(filter[0]))
			newImgData.Filter(iwsimage.RedFilterGenerator(filter[1]))
			newImgData.Filter(iwsimage.BlueFilterGenerator(filter[2]))
			
			newImgData.SaveFile(tmpDest)
			imageChan <- newImgData
		}(colorCollection[i*3:i*3+3], i<3)
	}
	
	imgData1 := <- imageChan
	
	imgData1.AssembleLeft(<- imageChan)
	imgData2 := <- imageChan
	imgData2.AssembleLeft(<- imageChan)

	imgData1.AssembleTop(imgData2)

	if err := imgData1.SaveFile(dest); err != nil {
		panic(err.Error())
	}
}
