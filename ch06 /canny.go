package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/building.jpg", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("image read fail")
		os.Exit(-1)
	}
	defer src.Close()

	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	dst := gocv.NewMat()
	defer dst.Close()

	gocv.Canny(src, &dst, 50, 150)
	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()
	dstWin.IMShow(dst)

	gocv.WaitKey(0)
}
