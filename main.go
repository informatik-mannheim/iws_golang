package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

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
