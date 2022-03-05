package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/namecard.jpg", gocv.IMReadAnyColor)
	if src.Empty() {
		fmt.Println("image read error!!")
		os.Exit(1)
	}
	defer src.Close()
	srcW := gocv.NewWindow("src")
	defer srcW.Close()
	srcW.IMShow(src)

	w, h := 720, 400
	srcQuad := gocv.NewPointVectorFromPoints([]image.Point{{325, 307}, {760, 369}, {718, 611}, {231, 515}})
	dstQuad := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {w - 1, 0}, {w - 1, h - 1}, {0, h - 1}})
	pers := gocv.GetPerspectiveTransform(srcQuad, dstQuad)

	dst := gocv.NewMat()
	defer dst.Close()
	gocv.WarpPerspective(src, &dst, pers, image.Pt(w, h))

	dstW := gocv.NewWindow("dst")
	defer dstW.Close()
	dstW.IMShow(dst)

	gocv.WaitKey(0)
}
