package main

import (
	"gocv.io/x/gocv"
	"image/color"
	"math/rand"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/milkdrop.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		println("image load fail!!")
		os.Exit(-1)
	}
	defer src.Close()
	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	srcBin := gocv.NewMat()
	defer srcBin.Close()
	gocv.Threshold(src, &srcBin, 0, 255, gocv.ThresholdOtsu)
	srcBinWin := gocv.NewWindow("src_bin")
	defer srcBinWin.Close()
	srcBinWin.IMShow(srcBin)

	contours := gocv.FindContours(srcBin, gocv.RetrievalList, gocv.ChainApproxNone)
	dst := gocv.Zeros(src.Rows(), src.Cols(), gocv.MatTypeCV8U)
	gocv.CvtColor(dst, &dst, gocv.ColorGrayToBGR)
	//dst := gocv.NewMatWithSizes([]int{src.Rows(), src.Cols()}, gocv.MatTypeCV8U)
	defer dst.Close()

	for i := 0; i < contours.Size(); i++ {
		c := color.RGBA{B: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), R: uint8(rand.Intn(255))}
		gocv.DrawContours(&dst, contours, i, c, 1)
	}

	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()
	dstWin.IMShow(dst)

	gocv.WaitKey(0)
}
