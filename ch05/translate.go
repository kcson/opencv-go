package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"os"
)

type Point func(x, y int) image.Point

func pt(x, y int) image.Point {
	return image.Point{X: x, Y: y}
}

func main() {
	src := gocv.IMRead("imgs/tekapo.bmp", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("image read fail!!")
		os.Exit(1)
	}
	defer src.Close()
	srcW := gocv.NewWindow("src")
	defer srcW.Close()
	srcW.IMShow(src)

	dst := gocv.NewMat()
	defer dst.Close()

	srcP := gocv.NewPoint2fVectorFromPoints([]gocv.Point2f{{0, 0}, {0, float32(src.Rows())}, {float32(src.Cols()), 0}})
	dstP := gocv.NewPoint2fVectorFromPoints([]gocv.Point2f{{200, 100}, {200, float32(src.Rows()) + 100}, {float32(src.Cols()) + 200, 100}})
	m := gocv.GetAffineTransform2f(srcP, dstP)
	//m := gocv.NewMatWithSize(2, 3, gocv.MatTypeCV32F)
	//m.SetFloatAt(0,0,1)
	//m.SetFloatAt(0,1,0)
	//m.SetFloatAt(0,2,200)
	//
	//m.SetFloatAt(1,0,0)
	//m.SetFloatAt(1,1,1)
	//m.SetFloatAt(1,2,100)

	gocv.WarpAffine(src, &dst, m, pt(0, 0))

	dstW := gocv.NewWindow("dst")
	defer dstW.Close()
	dstW.IMShow(dst)

	gocv.WaitKey(0)
}
