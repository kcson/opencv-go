package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	ref := gocv.IMRead("imgs/kids1.png", gocv.IMReadColor)
	mask := gocv.IMRead("imgs/kids1_mask.bmp", gocv.IMReadGrayScale)
	if ref.Empty() || mask.Empty() {
		fmt.Println("img read error!")
		os.Exit(1)
	}
	defer ref.Close()
	defer mask.Close()

	refYcrcb := gocv.NewMat()
	defer refYcrcb.Close()
	gocv.CvtColor(ref, &refYcrcb, gocv.ColorBGRToYCrCb)

	hist := gocv.NewMat()
	defer hist.Close()
	gocv.CalcHist([]gocv.Mat{refYcrcb}, []int{1, 2}, mask, &hist, []int{128, 128}, []float64{0, 256, 0, 256}, false)

	src := gocv.IMRead("imgs/kids2.png", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("src img read fail!!")

	}
	srcYcrcb := gocv.NewMat()
	defer srcYcrcb.Close()
	gocv.CvtColor(src, &srcYcrcb, gocv.ColorBGRToYCrCb)

	backProj := gocv.NewMat()
	defer backProj.Close()

	gocv.CalcBackProject([]gocv.Mat{srcYcrcb}, []int{1, 2}, hist, &backProj, []float64{0, 256, 0, 256}, true)

	srcWindow := gocv.NewWindow("src")
	defer srcWindow.Close()
	backProjWindow := gocv.NewWindow("backproj")
	defer backProjWindow.Close()

	srcWindow.IMShow(src)
	backProjWindow.IMShow(backProj)

	gocv.WaitKey(0)
}
