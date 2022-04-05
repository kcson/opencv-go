package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"opencv-go/common"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/lenna.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("image read fail!!")
		os.Exit(-1)
	}
	defer src.Close()

	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	kernel := common.NewMat(
		3,
		3,
		[][]float32{
			{-1, 0, 1},
			{-2, 0, 2},
			{-1, 0, 1},
		}, gocv.MatTypeCV32F)
	defer kernel.Close()

	dx := gocv.NewMat()
	defer dx.Close()
	dy := gocv.NewMat()
	defer dy.Close()

	//gocv.Filter2D(src, &dx, -1, *kernel, image.Pt(-1, -1), 128, gocv.BorderDefault)

	//sobel 함수
	gocv.Sobel(src, &dx, gocv.MatTypeCV32F, 1, 0, 3, 1, 0, gocv.BorderDefault)
	gocv.Sobel(src, &dy, gocv.MatTypeCV32F, 0, 1, 3, 1, 0, gocv.BorderDefault)

	mag := gocv.NewMat()
	defer mag.Close()

	gocv.Magnitude(dx, dy, &mag)
	mag.ConvertTo(&mag, gocv.MatTypeCV8U)
	magWin := gocv.NewWindow("mag")
	defer magWin.Close()
	magWin.IMShow(mag)

	edge := gocv.Zeros(mag.Rows(), mag.Cols(),  gocv.MatTypeCV8U)
	defer edge.Close()

	//gocv.Threshold(mag, &edge, 80, 255, gocv.ThresholdBinary)
	var threshold uint8 = 120
	for i := 0; i < mag.Rows(); i++ {
		for j := 0; j < mag.Cols(); j++ {
			if mag.GetUCharAt(i, j) > threshold {
				edge.SetUCharAt(i, j, 255)
			}
		}
	}

	edgeWin := gocv.NewWindow("edge")
	defer edgeWin.Close()

	edgeWin.IMShow(edge)

	gocv.WaitKey(0)
}
