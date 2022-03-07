package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"os"
)

func main() {
	capture, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println("Video capture error !!")
		os.Exit(1)
	}
	defer capture.Close()

	src := gocv.NewMat()
	defer src.Close()
	srcW := gocv.NewWindow("src")
	defer srcW.Close()
	camMode := 0
	for {
		if !capture.Read(&src) {
			break
		}
		if camMode == 1 {
			src = cartoonFilter(src)
		} else if camMode == 2 {
			src = pencilSketch(src)
			gocv.CvtColor(src, &src, gocv.ColorGrayToBGR)
		}
		srcW.IMShow(src)
		key := gocv.WaitKey(1)
		if key == 27 {
			break
		} else if key == 32 {
			camMode += 1
			if camMode == 3 {
				camMode = 0
			}
		}
	}

}

func cartoonFilter(src gocv.Mat) gocv.Mat {
	h, w := src.Rows(), src.Cols()

	img2 := gocv.NewMat()
	defer img2.Close()
	gocv.Resize(src, &img2, image.Point{X: w / 2, Y: h / 2}, 0, 0, gocv.InterpolationNearestNeighbor)

	blr := gocv.NewMat()
	defer blr.Close()
	gocv.BilateralFilter(img2, &blr, -1, 10, 7)

	edge := gocv.NewMat()
	defer edge.Close()
	gocv.Canny(img2, &edge, 80, 120)
	edge.ConvertToWithParams(&edge, gocv.MatTypeCV8U, -1, 255)
	gocv.CvtColor(edge, &edge, gocv.ColorGrayToBGR)

	dst := gocv.NewMat()
	gocv.BitwiseAnd(blr, edge, &dst)
	gocv.Resize(dst, &dst, image.Point{X: w, Y: h}, 0, 0, gocv.InterpolationNearestNeighbor)

	return dst
}

func pencilSketch(src gocv.Mat) gocv.Mat {
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)
	gray.ConvertTo(&gray, gocv.MatTypeCV32F)

	blr := gocv.NewMat()
	defer blr.Close()
	gocv.GaussianBlur(gray, &blr, image.Point{X: 0, Y: 0}, 3, 0, gocv.BorderDefault)
	gocv.Pow(blr, -1.0, &blr)

	dst := gocv.NewMat()
	gocv.MultiplyWithParams(gray, blr, &dst, 255, gray.Type())
	dst.ConvertTo(&dst, gocv.MatTypeCV8U)

	return dst
}
