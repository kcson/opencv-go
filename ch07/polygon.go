package main

import (
	"gocv.io/x/gocv"
	"image/color"
	"math"
	"os"
)

func main() {
	img := gocv.IMRead("imgs/polygon.bmp", gocv.IMReadColor)
	if img.Empty() {
		println("image read fail!!")
		os.Exit(1)
	}
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	srcBin := gocv.NewMat()
	defer srcBin.Close()
	gocv.Threshold(gray, &srcBin, 0, 255, gocv.ThresholdBinaryInv|gocv.ThresholdOtsu)
	contours := gocv.FindContours(srcBin, gocv.RetrievalExternal, gocv.ChainApproxNone)
	for i := 0; i < contours.Size(); i++ {
		pts := contours.At(i)
		if gocv.ContourArea(pts) < 400 {
			continue
		}
		approx := gocv.ApproxPolyDP(pts, gocv.ArcLength(pts, true)*0.02, true)
		vtc := approx.Size()
		if vtc == 3 {
			setLabel(&img, pts, "TRI")
		} else if vtc == 4 {
			setLabel(&img, pts, "RECT")
		} else {
			length := gocv.ArcLength(pts, true)
			area := gocv.ContourArea(pts)
			ratio := 4.0 * math.Pi * area / (length * length)
			if ratio > 0.85 {
				setLabel(&img, pts, "CIR")
			}
		}
	}
	imgWin := gocv.NewWindow("img")
	defer imgWin.Close()
	imgWin.IMShow(img)

	gocv.WaitKey(0)
}

func setLabel(img *gocv.Mat, pts gocv.PointVector, label string) {
	rect := gocv.BoundingRect(pts)
	gocv.Rectangle(img, rect, color.RGBA{R: 255}, 1)
	gocv.PutText(img, label, rect.Min, gocv.FontHersheyPlain, 1, color.RGBA{R: 255}, 1)
}
