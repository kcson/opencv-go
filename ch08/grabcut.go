package main

import (
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/nemo.jpg", gocv.IMReadColor)
	if src.Empty() {
		println("image read fail!!")
		os.Exit(1)
	}
	defer src.Close()
	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	roiWin := gocv.NewWindow("roi")
	defer roiWin.Close()
	rc := roiWin.SelectROI(src)

	mask := gocv.Zeros(src.Rows(), src.Cols(), gocv.MatTypeCV8U)
	defer mask.Close()

	bgModel := gocv.NewMat()
	defer bgModel.Close()
	fgModel :=gocv.NewMat()
	defer fgModel.Close()

	gocv.GrabCut(src, &mask, rc, &bgModel, &fgModel, 5, gocv.GCInitWithRect)

	mask.MultiplyUChar(64)

	maskWin := gocv.NewWindow("mask")
	defer maskWin.Close()
	maskWin.IMShow(mask)

	gocv.WaitKey(0)
}
