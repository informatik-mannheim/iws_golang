package iwsimage

import (
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
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
	if err := imgData1.LoadFile(src); err != nil {
		panic(err.Error())
	}

	imgData2 := imgData1.Copy()
	imgData3 := imgData1.Copy()
	imgData4 := imgData1.Copy()

	imgData1.Filter(GreenFilterGenerator(1.2))
	imgData1.Filter(RedFilterGenerator(0.7))
	imgData1.Filter(BlueFilterGenerator(0.7))

	dest1 := filepath.Join(filepath.Dir(dest), "img_1.bmp")
	if err := imgData1.SaveFile(dest1); err != nil {
		panic(err.Error())
	}

	imgData2.Filter(GreenFilterGenerator(0.7))
	imgData2.Filter(RedFilterGenerator(1.8))
	imgData2.Filter(BlueFilterGenerator(0.7))
	
	dest2 := filepath.Join(filepath.Dir(dest), "img_2.bmp")
	if err := imgData2.SaveFile(dest2); err != nil {
		panic(err.Error())
	}


	imgData3.Filter(GreenFilterGenerator(0.7))
	imgData3.Filter(RedFilterGenerator(0.7))
	imgData3.Filter(BlueFilterGenerator(1.8))
	
	dest3 := filepath.Join(filepath.Dir(dest), "img_3.bmp")
	if err := imgData3.SaveFile(dest3); err != nil {
		panic(err.Error())
	}
	
	imgData4.Filter(GreenFilterGenerator(0.3))
	imgData4.Filter(RedFilterGenerator(1.6))
	imgData4.Filter(BlueFilterGenerator(1.6))

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
