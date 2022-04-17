package main

import (
	"gocv.io/x/gocv"
	"image/color"
	"math/rand"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/contours.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		println("image read fail!!")
		os.Exit(-1)
	}
	defer src.Close()
	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	hier := gocv.NewMat()
	defer hier.Close()
	v := gocv.FindContoursWithParams(src, &hier, gocv.RetrievalCComp, gocv.ChainApproxNone)

	dst := gocv.NewMat()
	defer dst.Close()
	gocv.CvtColor(src, &dst, gocv.ColorGrayToBGR)
	idx := 0
	for {
		c := color.RGBA{B: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)),R: uint8(rand.Intn(255))}
		//gocv.DrawContours(&dst, v, idx, color.RGBA{B: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), R: uint8(rand.Intn(255))}, 3)
		drawContoursWithHierarchy(&dst, v, idx, c, 3, hier)
		vect := hier.GetVeciAt(0, idx)
		idx = int(vect[0])
		println(idx)
		if idx < 0 {
			break
		}
	}
	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()
	dstWin.IMShow(dst)

	gocv.WaitKey(0)
}

func drawContoursWithHierarchy(img *gocv.Mat, contours gocv.PointsVector, contourIdx int, c color.RGBA, thickness int, hier gocv.Mat) {
	gocv.DrawContours(img, contours, contourIdx, c, thickness)
	vect := hier.GetVeciAt(0, contourIdx)
	childIdx := int(vect[2])
	if childIdx >= 0 {
		drawContoursWithHierarchy(img, contours, childIdx, c, thickness, hier)
	}
}
