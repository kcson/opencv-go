package main

import (
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/sudoku.jpg", gocv.IMReadGrayScale)
	if src.Empty() {
		println("Image load fail!!")
		os.Exit(1)
	}
	defer src.Close()

	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()

	bs := 0
	dst := gocv.NewMat()
	defer dst.Close()
	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()
	dstTrackBar := dstWin.CreateTrackbarWithValue("Block size", &bs, 200)
	dstTrackBar.SetPos(11)

	i := 0
	for {
		i++
		if bs%2 == 0 {
			bs = bs-1
		}
		if bs < 3 {
			bs = 3
		}
		println(bs)
		gocv.AdaptiveThreshold(src, &dst, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, bs, 5)
		dstWin.IMShow(dst)
		if gocv.WaitKey(1000) == 27 {
			break
		}
	}
}
