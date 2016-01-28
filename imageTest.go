package main

import (
	"github.com/informatik-mannheim/iws_golang/iwsimage"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
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
		optimizedColorCollage(path.Join(currentDir, srcPicturePath), path.Join(usr.HomeDir, outputPicturePath), colorCollection)
	}

	// print the time difference between the runs
	timeNew := time.Since(timeNow).Nanoseconds() / int64(numberOfRuns)

	log.Println("average runtime(ns) -optimized:")
	log.Println(timeNew)
	log.Println("average runtime(ns) -not optimized:")
	log.Println(timeOld)

	log.Printf("Difference: %v ns\n", (timeOld - timeNew))
}

func optimizedColorCollage(src, dest string, colorCollection [4][3]float64) {

	imgData := iwsimage.NewImageData()
	if err := imgData.LoadFile(src); err != nil {
		panic(err.Error())
	}

	imageCopyChan := make(chan *iwsimage.ImageData, 4)
    imageChan := make(chan *iwsimage.ImageData, 4)
	workCounter := sync.WaitGroup{}
	imageCopyChan <- imgData
	for i := 0; i < 3; i++ {
		workCounter.Add(1)
		go func() {
			defer workCounter.Done()
			imageCopyChan <- imgData.Copy()
		}()
	}

	workCounter.Wait()
    chanLenght := len(imageCopyChan)
    
	for i := 0; i < chanLenght; i++ {
		tmpDest := filepath.Join(filepath.Dir(dest), "img_"+strconv.Itoa(i+1)+"_optimized.bmp")
		go func(filter [3]float64, newImgData *iwsimage.ImageData) {
			newImgData.Filter(iwsimage.OldGreenFilterGenerator(filter[0]))
			newImgData.Filter(iwsimage.OldRedFilterGenerator(filter[1]))
			newImgData.Filter(iwsimage.OldBlueFilterGenerator(filter[2]))

			if err := newImgData.SaveFile(tmpDest); err != nil {
				panic(err.Error())
			}
            
			imageChan <- newImgData
		}(colorCollection[i], <-imageCopyChan)
	}

	assembledImage := <-imageChan
	assembledImage.AssembleLeft(<-imageChan)
	imgDataTop := <-imageChan
	imgDataTop.AssembleLeft(<-imageChan)
    
    assembledImage.AssembleTop(imgDataTop)

	// Save assembled image
	if err := assembledImage.SaveFile(dest); err != nil {
		panic(err.Error())
	}
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
