package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image/color"
	"math"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/tekapo.bmp", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("image read fail!!")
		os.Exit(1)
	}
	defer src.Close()

	map1 := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV32FC1)
	defer map1.Close()
	map2 := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV32FC1)
	defer map2.Close()
	for x := 0; x < src.Cols(); x++ {
		for y := 0; y < src.Rows(); y++ {
			map1.SetFloatAt(y, x, float32(x))
			map2.SetFloatAt(y, x, float32(y)+10*float32(math.Sin(float64(x)/32)))
		}
	}
	dst := gocv.NewMat()
	defer dst.Close()

	gocv.Remap(src, &dst, &map1, &map2, gocv.InterpolationCubic, gocv.BorderDefault, color.RGBA{B: 0, G: 0, R: 0})
	dstW := gocv.NewWindow("dst")
	defer dstW.Close()

	dstW.IMShow(dst)

	gocv.WaitKey(0)
}
