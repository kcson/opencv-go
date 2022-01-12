package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"opencv-go/common"
	"os"
)

func main() {
	src := gocv.IMRead("Hawkes.jpg", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("Img read error!!")
		os.Exit(-1)
	}
	defer src.Close()
	dst := gocv.NewMat()
	defer dst.Close()

	gocv.EqualizeHist(src, &dst)

	srcW := gocv.NewWindow("src")
	defer srcW.Close()
	dstW := gocv.NewWindow("dst")
	defer dstW.Close()

	srcW.IMShow(src)
	dstW.IMShow(dst)

	srcHist := gocv.NewMat()
	defer srcHist.Close()
	dstHist := gocv.NewMat()
	defer dstHist.Close()

	gocv.CalcHist([]gocv.Mat{src}, []int{0}, gocv.NewMat(), &srcHist, []int{256}, []float64{0, 256}, false)
	gocv.CalcHist([]gocv.Mat{dst}, []int{0}, gocv.NewMat(), &dstHist, []int{256}, []float64{0, 256}, false)

	gocv.NewWindow("srcHist").IMShow(common.DrawHistogramV2(srcHist))
	gocv.NewWindow("dstHist").IMShow(common.DrawHistogramV2(dstHist))

	gocv.WaitKey(0)
}


