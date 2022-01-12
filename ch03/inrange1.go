package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	//src := gocv.IMRead("imgs/candies.png", gocv.IMReadColor)
	src := gocv.IMRead("imgs/candies2.png",gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("img read fail!!")
		os.Exit(1)
	}
	defer src.Close()
	srcHsv := gocv.NewMat()
	defer srcHsv.Close()

	dst1 := gocv.NewMat()
	defer dst1.Close()

	dst2 := gocv.NewMat()
	defer dst2.Close()

	gocv.CvtColor(src, &srcHsv, gocv.ColorBGRToHSV)

	gocv.InRangeWithScalar(src, gocv.Scalar{Val1: 0, Val2: 128, Val3: 0}, gocv.Scalar{Val1: 100, Val2: 255, Val3: 100}, &dst1)
	gocv.InRangeWithScalar(srcHsv, gocv.Scalar{Val1: 90, Val2: 150, Val3: 0}, gocv.Scalar{Val1: 130, Val2: 255, Val3: 255}, &dst2)

	sw := gocv.NewWindow("src")
	defer sw.Close()

	d1 := gocv.NewWindow("dst1")
	defer d1.Close()

	d2 := gocv.NewWindow("dst2")
	defer d2.Close()

	sw.IMShow(src)
	d1.IMShow(dst1)
	d2.IMShow(dst2)

	gocv.WaitKey(0)

}
