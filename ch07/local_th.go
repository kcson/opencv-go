package main

import (
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/sudoku.jpg", gocv.IMReadGrayScale)
	if src.Empty() {
		println("Image load fail!!")
		os.Exit(1)
	}
	defer src.Close()
	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	dst1 := gocv.NewMat()
	defer dst1.Close()

	// 전역 이진화
	gocv.Threshold(src, &dst1, 0, 255, gocv.ThresholdOtsu)
	dst1Win := gocv.NewWindow("dst1")
	defer dst1Win.Close()
	dst1Win.IMShow(dst1)

	// 지역 이진화
	dst2 := gocv.Zeros(src.Rows(), src.Cols(), gocv.MatTypeCV8U)
	defer dst2.Close()
	bw := src.Cols() / 4
	bh := src.Rows() / 4

	for i := 0; i < 4; i++ {
		src_ := src.RowRange(i*bh, (i+1)*bh)
		dst_ := dst2.RowRange(i*bh, (i+1)*bh)
		for j := 0; j < 4; j++ {
			src__ := src_.ColRange(j*bw, (j+1)*bw)
			dst__ := dst_.ColRange(j*bw, (j+1)*bw)
			gocv.Threshold(src__, &dst__, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)
		}
	}

	dst2Win := gocv.NewWindow("dst2")
	defer dst2Win.Close()
	dst2Win.IMShow(dst2)

	gocv.WaitKey(0)
}
