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
		fmt.Println("image read fail!!")
		os.Exit(1)
	}
	defer src.Close()
	srcW := gocv.NewWindow("src")
	defer srcW.Close()
	srcW.IMShow(src)

	srcYcrcb := gocv.NewMat()
	defer srcYcrcb.Close()
	gocv.CvtColor(src, &srcYcrcb, gocv.ColorBGRToYCrCb)

	srcF := gocv.NewMat()
	defer srcF.Close()
	gocv.ExtractChannel(srcYcrcb, &srcF, 0)

	srcFW := gocv.NewWindow("srcFW")
	defer srcFW.Close()
	srcFW.IMShow(srcF)

	blr := gocv.NewMat()
	defer blr.Close()
	gocv.GaussianBlur(srcF, &blr, image.Point{X: 0, Y: 0}, 2.0, 0, gocv.BorderDefault)

	blrW := gocv.NewWindow("blr")
	defer blrW.Close()
	blrW.IMShow(blr)

	dst := gocv.NewMat()
	defer dst.Close()
	//gocv.AddWeighted(src, 1, blr, -1, 128, &dst)
	//2f - f^
	//gocv.AddWeighted(src, 2, blr, -1, 0, &dst)
	//gocv.AddWeighted(src, 2, blr, -1, 0, &dst)
	gocv.AddWeighted(srcF, 2, blr, -1, 0, &dst)
	//src.MultiplyFloat(2)
	//blr.MultiplyFloat(-1)
	//gocv.Add(src,blr,&dst)
	gocv.InsertChannel(dst, &srcYcrcb, 0)
	gocv.CvtColor(srcYcrcb, &dst, gocv.ColorYCrCbToBGR)

	dstW := gocv.NewWindow("dst")
	defer dstW.Close()
	dstW.IMShow(dst)

	gocv.WaitKey(0)
}
