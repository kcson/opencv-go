package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/cells.png", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("image read fail!!")
		os.Exit(1)
	}
	defer src.Close()
	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	dst := gocv.NewMat()
	defer dst.Close()
	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()

	threshold := 0
	trackBar := dstWin.CreateTrackbarWithValue("Threshold", &threshold, 255)
	trackBar.SetPos(150)

	for {
		gocv.Threshold(src, &dst, float32(threshold), 255, gocv.ThresholdBinary)
		dstWin.IMShow(dst)
		if gocv.WaitKey(3) == 27 {
			break
		}
	}
}
