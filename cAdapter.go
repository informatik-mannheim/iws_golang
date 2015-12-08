package iws_golang

import (
	"errors"
	"fmt"
	"unsafe"
)

/*
#include <stdio.h>
extern int c_reader(void* callback, char *src);
extern int c_writeDataToImageFile(int *data, char *fileName, int width, int height);
*/
import "C"
//export handelPixelData
func handelPixelData(object unsafe.Pointer, r, g, b int) {
	info := (*ImageData)(object)
	info.addPixel(r, g, b)
}

//export getPixelInfo
func getPixelInfo(object unsafe.Pointer, width, height int) {
	info := (*ImageData)(object)
	info.SetSize(width, height)
}

// loadFile loads a bitmap file via c
func loadFile(src string, imgData *ImageData) error {

	c := make(chan int)
	go (func() {
		defer func() {
			if r := recover(); r != nil {
				c <- 2
			}
		}()
		errorCode := int(C.c_reader(unsafe.Pointer(imgData), C.CString(src)))
		c <- errorCode

	})()

	errorCode := <-c

	if errorCode == 1 {
		return errors.New("File not found")
	} else if errorCode == 2 {
		return errors.New("No Bitmap")
	}

	return nil

}

// savefile saves a bitmap file via c
func savefile(destGo string, imgData *ImageData) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	dest := C.CString(destGo)

	cCastInt := make([]C.int, 0, imgData.Height*imgData.Height*3)
	for _, i := range imgData.PixelData {
		cCastInt = append(cCastInt, C.int(i))
	}

	errorCode := int(C.c_writeDataToImageFile(&cCastInt[0], dest, C.int(imgData.Width), C.int(imgData.Height)))

	if errorCode > 0 {
		err = errors.New("Not able to save File")
	}

	return
}
