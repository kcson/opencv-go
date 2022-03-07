package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/rose.bmp", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("image read fail!!!")
		os.Exit(1)
	}
	defer src.Close()

	srcW := gocv.NewWindow("src")
	defer srcW.Close()
	srcW.IMShow(src)

	dst := gocv.NewMat()
	defer dst.Close()
	gocv.Resize(src, &dst, image.Pt(0, 0), 4, 4, gocv.InterpolationLanczos4)

	dstW := gocv.NewWindow("dst")
	defer dstW.Close()
	dstW.IMShow(dst)

	gocv.WaitKey(0)
}
