package main

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/keyboard.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		println("Image load fail!!")
		os.Exit(1)
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

	labels := gocv.NewMat()
	defer labels.Close()
	stats := gocv.NewMat()
	defer stats.Close()
	centroids := gocv.NewMat()
	defer centroids.Close()
	cnt := gocv.ConnectedComponentsWithStats(srcBin, &labels, &stats, &centroids)

	dst := gocv.NewMat()
	defer dst.Close()
	gocv.CvtColor(src, &dst, gocv.ColorGrayToBGR)

	for i := 1; i < cnt; i++ {
		x, y, w, h, area := stats.GetIntAt(i, 0), stats.GetIntAt(i, 1), stats.GetIntAt(i, 2), stats.GetIntAt(i, 3), stats.GetIntAt(i, 4)
		if area < 20 {
			continue
		}
		gocv.Rectangle(&dst, image.Rect(int(x), int(y), int(x+w), int(y+h)), color.RGBA{R: 255}, 1)
	}
	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()
	dstWin.IMShow(dst)

	gocv.WaitKey(0)
}
