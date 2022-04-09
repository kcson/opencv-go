package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/sudoku.jpg", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("Image load fail!!")
		os.Exit(1)
	}
	defer src.Close()

	dst := gocv.NewMat()
	defer dst.Close()
	th := gocv.Threshold(src, &dst, 0, 255, gocv.ThresholdOtsu)
	println("threshold : ", th)

	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()

	srcWin.IMShow(src)
	dstWin.IMShow(dst)

	gocv.WaitKey(0)
}
