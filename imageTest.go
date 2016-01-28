package main

import (
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"time"
	//	"sync"
	"github.com/informatik-mannheim/iws_golang/iwsimage"
)

// This is the main entry point for an image manipulation program
func runImageTest() {
	// defines how often the image manipulation is done
	const numberOfRuns = 1
	const srcPicturePath = "pictures/bridge.bmp"
	const outputPicturePath = "Desktop/bridgecolored.bmp"

	colorCollection := [...][3]float64{{1.2, 0.7, 0.7}, {0.7, 1.8, 0.7}, {0.7, 0.7, 1.8}, {0.3, 1.6, 1.6}}

	// finds current directory Path
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// finds current users for home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	timeNow := time.Now()

	for i := 0; i < numberOfRuns; i++ {
		//execute image manipulation with old algorithm
		oldColorCollage(path.Join(currentDir, srcPicturePath), path.Join(usr.HomeDir, outputPicturePath), colorCollection)
	}

	timeOld := time.Since(timeNow).Nanoseconds() / int64(numberOfRuns)

	timeNow = time.Now()
	for i := 0; i < numberOfRuns; i++ {
		//execute image manipulation again with with optimized algorithm
		// Optimize this version by using better goroutines
//		colorCollage(path.Join(currentDir, srcPicturePath), path.Join(usr.HomeDir, outputPicturePath), colorCollection)
	}

	// print the time difference between the runs
	timeNew := time.Since(timeNow).Nanoseconds() / int64(numberOfRuns)

	log.Println("average runtime(ns) -optimized:")
	log.Println(timeNew)
	log.Println("average runtime(ns) -not optimized:")
	log.Println(timeOld)

	log.Printf("Difference: %v ns\n", (timeOld - timeNew))
}

// optimized version
func colorCollage(src, dest string, collerCollection [4][3]float64) {

	imgData := iwsimage.NewImageData()
	if err := imgData.LoadFile(src); err != nil {
		panic(err.Error())
	}

	imageChan := make(chan *iwsimage.ImageData, 4)
	//	workCounter := sync.WaitGroup{}
	//	imageChan <- imgData
	//	for i := 0; i < 3; i++ {
	//		workCounter.Add(1)
	//		go func() {
	//			defer workCounter.Done()
	//			imageChan <- imgData.Copy()
	//		}()
	//	}

	//	workCounter.Wait();

	for i := 0; i < 4; i++ {
		//		tmpDest := filepath.Join(filepath.Dir(dest), "img_" + strconv.Itoa(i+1) + "_p.bmp")
		//		go func(filter []float64,  newImgData *iwsimage.ImageData) {
		//			newImgData.Filter(iwsimage.GreenFilterGenerator(filter[0]))
		//			newImgData.Filter(iwsimage.RedFilterGenerator(filter[1]))
		//			newImgData.Filter(iwsimage.BlueFilterGenerator(filter[2]))

		//			newImgData.SaveFile(tmpDest)
		//			imageChan <- newImgData
		//		}(colorCollection[i*3:i*3+3], <-imageChan)
	}

	imgData1 := <-imageChan
	//
	imgData1.AssembleLeft(<-imageChan)
	//	imgData2 := <- imageChan
	//	imgData2.AssembleLeft(<- imageChan)

	//	imgData1.AssembleTop(imgData2)

	if err := imgData1.SaveFile(dest); err != nil {
		panic(err.Error())
	}
}

// Gets a path to a source and destination bitmapfile
// Creates four copies of the file and manipulates the color of each copy
// Assembles the four changed pictures together and writes the final picture to dest
func oldColorCollage(src, dest string, colorCollection [4][3]float64) {
	imgData := iwsimage.NewImageData()
	if err := imgData.LoadFile(src); err != nil {
		panic(err.Error())
	}

	imageList := [...]*iwsimage.ImageData{
		imgData.Copy(),
		imgData.Copy(),
		imgData.Copy(),
		imgData}

	for index, newImgData := range imageList {
		newImgData.Filter(iwsimage.OldGreenFilterGenerator(colorCollection[index][0]))
		newImgData.Filter(iwsimage.OldRedFilterGenerator(colorCollection[index][1]))
		newImgData.Filter(iwsimage.OldBlueFilterGenerator(colorCollection[index][2]))
		dest := filepath.Join(filepath.Dir(dest), "coloredImg_"+strconv.Itoa(index+1)+".bmp")
		if err := newImgData.SaveFile(dest); err != nil {
			panic(err.Error())
		}
	}

	imageList[0].AssembleLeft(imageList[1])
	imageList[2].AssembleLeft(imageList[3])
	imageList[0].AssembleTop(imageList[2])

	if err := imageList[0].SaveFile(dest); err != nil {
		panic(err.Error())
	}
}
