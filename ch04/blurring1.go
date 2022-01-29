package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/rose.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("image read fail !!")
		os.Exit(1)
	}
	defer src.Close()

	srcWindow := gocv.NewWindow("src")
	defer srcWindow.Close()

	srcWindow.IMShow(src)

	dst := gocv.NewMat()
	defer dst.Close()

	filter := gocv.Ones(3, 3, gocv.MatTypeCV64F)
	filter.MultiplyFloat(1.0 / 9.0)
	defer filter.Close()

	dstWindow := gocv.NewWindow("dst")
	defer dstWindow.Close()

	//dstWindow.IMShow(dst)

	kSize := []int{3, 5, 7}
	//gocv.Filter2D(src, &dst, -1, filter, image.Point{X: -1, Y: -1}, 0, gocv.BorderDefault)
	//gocv.Blur(src, &dst, image.Point{X: 5, Y: 5})
	for _, size := range kSize {
		gocv.Blur(src, &dst, image.Point{X: size, Y: size})
		gocv.PutTextWithParams(
			&dst,
			fmt.Sprintf("Mean : %d X %d", size, size), image.Point{X: 10, Y: 30},
			gocv.FontHersheySimplex,
			1.0,
			color.RGBA{B: 255, G: 255, R: 255},
			2,
			gocv.LineAA,
			false)
		dstWindow.IMShow(dst)
		gocv.WaitKey(0)
	}

	//gocv.WaitKey(0)
}
