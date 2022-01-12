package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"opencv-go/common"
	"os"
)

func main() {
	src := gocv.IMRead("lenna.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("Img read fail!!")
		os.Exit(2)
	}

	hist := gocv.NewMat()
	gocv.CalcHist([]gocv.Mat{src}, []int{0}, gocv.NewMat(), &hist, []int{256}, []float64{0, 256}, false)

	gocv.NewWindow("src").IMShow(src)
	gocv.WaitKey(1)
	gocv.NewWindow("hist").IMShow(common.DrawHistogramV2(hist))
	gocv.WaitKey(0)
}

