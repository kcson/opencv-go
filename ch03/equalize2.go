package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/field.bmp", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("img read fail")
		os.Exit(-1)
	}
	defer src.Close()

	dst := gocv.NewMat()
	defer dst.Close()

	gocv.CvtColor(src, &dst, gocv.ColorBGRToYCrCb)

	planes := gocv.Split(dst)
	gocv.EqualizeHist(planes[0], &planes[0])
	gocv.Merge(planes,&dst)
	gocv.CvtColor(dst, &dst, gocv.ColorYCrCbToBGR)

	sw := gocv.NewWindow("src")
	defer sw.Close()
	dw := gocv.NewWindow("dst")
	defer dw.Close()

	sw.IMShow(src)
	dw.IMShow(dst)

	gocv.WaitKey(0)
}
