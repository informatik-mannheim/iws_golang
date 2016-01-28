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
	const (
		numberOfRuns      = 3
		srcPicturePath    = "pictures/bridge.bmp"
		outputPicturePath = "Desktop/assembled_Bridge_Image.bmp"
	)

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
		// execute image manipulation again with with optimized algorithm
		// Optimize this version by using better goroutines
		//		optimizedColorCollage(path.Join(currentDir, srcPicturePath), path.Join(usr.HomeDir, outputPicturePath), colorCollection)
	}

	// print the time difference between the runs
	timeNew := time.Since(timeNow).Nanoseconds() / int64(numberOfRuns)

	log.Println("average runtime(ns) -optimized:")
	log.Println(timeNew)
	log.Println("average runtime(ns) -not optimized:")
	log.Println(timeOld)

	log.Printf("Difference: %v ns\n", (timeOld - timeNew))
}

// modify here
func optimizedColorCollage(src, dest string, collerCollection [4][3]float64) {

	// 1. Ask yourself what parts of the algorithm in OldColorCollage could be faster concurrently
	// 2. Use goroutines to spawn threads
	// 3. Use Channels to communicate between goroutines
	// 4. Find a way to synchronize the output of goroutines. Hint you could use sync.WaitGroup{} for that

	imgData := iwsimage.NewImageData()
	if err := imgData.LoadFile(src); err != nil {
		panic(err.Error())
	}
	/**  Hints **/
	// use channels to synchronize different goroutines without explicit locks or condition variables
	// create a channel: make(chan datatype, buffersize)
	// fill channel: chan <- input
	// read out of channel: result := <- chan
	// see also: https://tour.golang.org/concurrency/2

	// channel with pointers on ImageData
	// imageChan := make(chan *iwsimage.ImageData, 4)

	// We need 4 copies of the image: see func oldColorCollage for the operations

	// We need three filter operations for each image
	// And afterwards the image needs to be saved
	// save all 4 images to desktop with:
	// tmpDest := filepath.Join(filepath.Dir(dest), "img_" + strconv.Itoa(i+1) + "_modified.bmp")
	// imgData.SaveFile(tmpDest)

	// 	Assemble the images at the end as done in func oldColorCollage
	//	imgData.AssembleLeft(<-imageChan)

	// Save assembled image
	//	if err := assembledImage.SaveFile(dest); err != nil {
	//		panic(err.Error())
	//	}
}

// Gets a path to a source and destination bitmapfile
// Creates four copies of the file and manipulates the color of each copy
// Assembles (=zusammenfügen) the four changed pictures together into one picture and writes the final picture to dest
func oldColorCollage(src, dest string, colorCollection [4][3]float64) {

	// read bitmap from given source
	imgData := iwsimage.NewImageData()
	if err := imgData.LoadFile(src); err != nil {
		panic(err.Error())
	}

	// create working copies of the given image
	imageList := [...]*iwsimage.ImageData{
		imgData.Copy(),
		imgData.Copy(),
		imgData.Copy(),
		imgData}

	// runs colorfilters on each one of the images
	for index, newImgData := range imageList {
		newImgData.Filter(iwsimage.OldGreenFilterGenerator(colorCollection[index][0]))
		newImgData.Filter(iwsimage.OldRedFilterGenerator(colorCollection[index][1]))
		newImgData.Filter(iwsimage.OldBlueFilterGenerator(colorCollection[index][2]))
		dest := filepath.Join(filepath.Dir(dest), "coloredImg_"+strconv.Itoa(index+1)+".bmp")
		// save the filtered image to the desktop
		if err := newImgData.SaveFile(dest); err != nil {
			panic(err.Error())
		}
	}

	// assembles the 4 images together into one final image  
	assembledImage := imageList[0]
	assembledImage.AssembleLeft(imageList[1])
	imageList[2].AssembleLeft(imageList[3])
	assembledImage.AssembleTop(imageList[2])

	// saves the assembled image to desktop
	if err := assembledImage.SaveFile(dest); err != nil {
		panic(err.Error())
	}
}
