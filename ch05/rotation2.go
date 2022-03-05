package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/tekapo.bmp", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("img read error!!")
		os.Exit(1)
	}
	defer src.Close()

	centerX, centerY := src.Cols()/2, src.Rows()/2

	rot := gocv.GetRotationMatrix2D(image.Pt(centerX, centerY), 20, 1)

	dst := gocv.NewMat()
	defer dst.Close()
	gocv.WarpAffine(src, &dst, rot, image.Pt(0, 0))

	dstW := gocv.NewWindow("dst")
	defer dstW.Close()
	dstW.IMShow(dst)

	gocv.WaitKey(0)
}
