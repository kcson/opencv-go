package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/lenna.bmp", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("image read fail !!")
		os.Exit(1)
	}
	defer src.Close()

	srcW := gocv.NewWindow("src")
	defer srcW.Close()
	srcW.IMShow(src)

	dst := gocv.NewMat()
	defer dst.Close()
	gocv.BilateralFilter(src, &dst, -1, 10, 5)

	dstW := gocv.NewWindow("dst")
	defer dstW.Close()
	dstW.IMShow(dst)

	gocv.WaitKey(0)
}
