package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/cropland.png", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("img read fail!!")
		os.Exit(1)
	}
	defer src.Close()
	srcWindow := gocv.NewWindow("src")
	defer srcWindow.Close()
	srcWindow.IMShow(src)

	roi := srcWindow.SelectROI(src)

	srcYcrcb := gocv.NewMat()
	defer srcYcrcb.Close()

	gocv.CvtColor(src, &srcYcrcb, gocv.ColorBGRToYCrCb)
	crop := srcYcrcb.ColRange(roi.Min.X, roi.Min.X+roi.Dx())
	crop = crop.RowRange(roi.Min.Y, roi.Min.Y+roi.Dy())

	hist := gocv.NewMat()
	defer hist.Close()
	gocv.CalcHist([]gocv.Mat{crop}, []int{1, 2}, gocv.NewMat(), &hist, []int{128, 128}, []float64{0, 256, 0, 256}, false)

	histWindow := gocv.NewWindow("hist")
	defer histWindow.Close()
	histWindow.IMShow(hist)

	backProj := gocv.NewMat()
	defer backProj.Close()
	gocv.CalcBackProject([]gocv.Mat{srcYcrcb}, []int{1, 2}, hist, &backProj, []float64{0, 256, 0, 256}, true)

	backProjWindow := gocv.NewWindow("backproj")
	defer backProjWindow.Close()
	backProjWindow.IMShow(backProj)

	dst := gocv.NewMat()
	defer dst.Close()

	src.CopyToWithMask(&dst,backProj)

	dstWindow := gocv.NewWindow("dst")
	defer dstWindow.Close()
	dstWindow.IMShow(dst)

	gocv.WaitKey(0)
}
