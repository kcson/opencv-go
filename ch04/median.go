package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/noise.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("image read error!!")
		os.Exit(1)
	}
	defer src.Close()
	srcW := gocv.NewWindow("src")
	defer srcW.Close()

	srcW.IMShow(src)

	dst := gocv.NewMat()
	defer dst.Close()
	gocv.MedianBlur(src, &dst, 3)
	dstW := gocv.NewWindow("dst")
	dstW.MoveWindow(src.Cols() + 10,0)
	defer dstW.Close()

	dstW.IMShow(dst)

	gocv.WaitKey(0)
}
