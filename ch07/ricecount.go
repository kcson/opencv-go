package main

import (
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/rice.png", gocv.IMReadGrayScale)
	if src.Empty() {
		println("image load fail!!")
		os.Exit(1)
	}
	defer src.Close()

	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	dst1 := gocv.Zeros(src.Rows(), src.Cols(), gocv.MatTypeCV8U)
	defer dst1.Close()
	labels := gocv.NewMat()
	defer labels.Close()

	bw := src.Cols() / 4
	bh := src.Rows() / 4
	for i := 0; i < 4; i++ {
		src_ := src.RowRange(i*bh, (i+1)*bh)
		dst_ := dst1.RowRange(i*bh, (i+1)*bh)
		for j := 0; j < 4; j++ {
			src__ := src_.ColRange(j*bw, (j+1)*bw)
			dst__ := dst_.ColRange(j*bw, (j+1)*bw)
			gocv.Threshold(src__, &dst__, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)
		}
	}
	dst1Win := gocv.NewWindow("dst1")
	defer dst1Win.Close()
	dst1Win.IMShow(dst1)

	//gocv.AdaptiveThreshold(src, &dst1, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdOtsu, 5, 5)
	cnt1 := gocv.ConnectedComponents(dst1, &labels)
	println("cnt1 : ", cnt1)

	dst2 := gocv.Zeros(src.Rows(), src.Cols(), gocv.MatTypeCV8U)
	defer dst2.Close()

	gocv.MorphologyEx(dst1, &dst2, gocv.MorphOpen, gocv.NewMat())
	cnt2 := gocv.ConnectedComponents(dst2, &labels)
	println("cnt2 : ", cnt2)

	dst2Win := gocv.NewWindow("dst2")
	defer dst2Win.Close()
	dst2Win.IMShow(dst2)

	gocv.WaitKey(0)
}
