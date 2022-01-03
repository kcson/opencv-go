package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("lenna.bmp", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("Image read error !!")
		os.Exit(2)
	}
	defer src.Close()

	window1 := gocv.NewWindow("src")
	defer window1.Close()

	window2 := gocv.NewWindow("dst")
	defer window2.Close()

	window1.IMShow(src)
	//window1.WaitKey(0)

	dst := gocv.NewMat()
	//src.AddUChar(100)
	gocv.Add(src, gocv.NewMatWithSizeFromScalar(gocv.Scalar{Val1: 100, Val2: 100, Val3: 100}, src.Rows(), src.Cols(), gocv.MatTypeCV8UC3), &dst)
	window2.IMShow(dst)
	window2.WaitKey(0)
}
