package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src1 := gocv.IMRead("lenna256.bmp", gocv.IMReadGrayScale)
	src2 := gocv.IMRead("square.bmp", gocv.IMReadGrayScale)

	if src1.Empty() || src2.Empty() {
		fmt.Println("Image load fail!!")
		os.Exit(0)
	}
	defer src1.Close()
	defer src2.Close()

	dst1 := gocv.NewMat()
	window1 := gocv.NewWindow("dst1")
	dst2 := gocv.NewMat()
	window2 := gocv.NewWindow("dst2")
	dst3 := gocv.NewMat()
	window3 := gocv.NewWindow("dst3")
	dst4 := gocv.NewMat()
	window4 := gocv.NewWindow("dst4")

	gocv.Add(src1, src2, &dst1)
	gocv.AddWeighted(src1, 0.5, src2, 0.5, 0.0, &dst2)
	gocv.Subtract(src1, src2, &dst3)
	gocv.AbsDiff(src1, src2, &dst4)

	window1.IMShow(dst1)
	window2.IMShow(dst2)
	window2.MoveWindow(dst1.Cols(), 0)
	window3.IMShow(dst3)
	window3.MoveWindow(dst1.Cols()+dst2.Cols(), 0)
	window4.IMShow(dst4)
	window4.MoveWindow(dst1.Cols()+dst2.Cols()+dst3.Cols(), 0)

	window1.WaitKey(0)
}
