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
	
	timeNew := time.Since(timeNow).Nanoseconds() / int64(numberOfRuns)
	
	timeNow = time.Now()
	for i := 0; i < numberOfRuns; i++ {
		oldColorCollage(path.Join(dir, "/test/pictures/bridge.bmp"), path.Join(usr.HomeDir, "Desktop/bridgecollor.bmp"))
	}
	
	timeOld := time.Since(timeNow).Nanoseconds() / int64(numberOfRuns)
	
	
	fmt.Println("average runtime(ns) -New:");
	fmt.Println(timeNew)
	fmt.Println("average runtime(ns) -Old:");
	fmt.Println(timeOld)
	
	fmt.Printf("Difference: %v ns\n", (timeOld - timeNew));
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

func oldColorCollage(src, dest string) {

	imgData1 := iwsimage.NewImageData()
	if err := imgData1.LoadFile(src); err != nil {
		panic(err.Error())
	}

	imgData2 := imgData1.Copy()
	imgData3 := imgData1.Copy()
	imgData4 := imgData1.Copy()

	imgData1.Filter(iwsimage.GreenFilterGenerator(1.2))
	imgData1.Filter(iwsimage.RedFilterGenerator(0.7))
	imgData1.Filter(iwsimage.BlueFilterGenerator(0.7))

	dest1 := filepath.Join(filepath.Dir(dest), "img_1.bmp")
	if err := imgData1.SaveFile(dest1); err != nil {
		panic(err.Error())
	}

	imgData2.Filter(iwsimage.GreenFilterGenerator(0.7))
	imgData2.Filter(iwsimage.RedFilterGenerator(1.8))
	imgData2.Filter(iwsimage.BlueFilterGenerator(0.7))

	dest2 := filepath.Join(filepath.Dir(dest), "img_2.bmp")
	if err := imgData2.SaveFile(dest2); err != nil {
		panic(err.Error())
	}

	imgData3.Filter(iwsimage.GreenFilterGenerator(0.7))
	imgData3.Filter(iwsimage.RedFilterGenerator(0.7))
	imgData3.Filter(iwsimage.BlueFilterGenerator(1.8))

	dest3 := filepath.Join(filepath.Dir(dest), "img_3.bmp")
	if err := imgData3.SaveFile(dest3); err != nil {
		panic(err.Error())
	}

	imgData4.Filter(iwsimage.GreenFilterGenerator(0.3))
	imgData4.Filter(iwsimage.RedFilterGenerator(1.6))
	imgData4.Filter(iwsimage.BlueFilterGenerator(1.6))

	dest4 := filepath.Join(filepath.Dir(dest), "img_4.bmp")
	if err := imgData4.SaveFile(dest4); err != nil {
		panic(err.Error())
	}

	imgData1.AssembleLeft(imgData2)
	imgData3.AssembleLeft(imgData4)

	imgData1.AssembleTop(imgData3)

	if err := imgData1.SaveFile(dest); err != nil {
		panic(err.Error())
	}
}