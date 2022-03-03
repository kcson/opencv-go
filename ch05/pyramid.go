package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/cat.bmp", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("image read fail !!")
		os.Exit(1)
	}
	defer src.Close()

	cpy := gocv.NewMat()
	defer cpy.Close()
	src.CopyTo(&cpy)
	gocv.Rectangle(&cpy, image.Rect(250, 120, 450, 320), color.RGBA{B: 0, G: 0, R: 255}, 2)
	cpyW := gocv.NewWindow("cpy")
	defer cpyW.Close()

	cpyW.IMShow(cpy)

	gocv.WaitKey(0)

	dst := gocv.NewMat()
	defer dst.Close()
	for i := 1; i <= 4; i++ {
		dstW := gocv.NewWindow("dst")
		gocv.PyrDown(src, &src, image.Pt(0, 0), gocv.BorderDefault)
		src.CopyTo(&cpy)
		gocv.RectangleWithParams(&cpy, image.Rect(250, 120, 450, 320), color.RGBA{B: 0, G: 0, R: 255}, 2, gocv.LineAA, i)
		dstW.IMShow(cpy)
		gocv.WaitKey(0)
		dstW.Close()
	}

}
