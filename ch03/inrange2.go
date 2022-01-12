package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/candies.png", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("image read fail!!")
		os.Exit(2)
	}
	defer src.Close()
	srcHsv := gocv.NewMat()
	defer srcHsv.Close()
	dst := gocv.NewMat()
	defer dst.Close()

	gocv.CvtColor(src, &srcHsv, gocv.ColorBGRToHSV)

	sw := gocv.NewWindow("src")
	defer sw.Close()

	sw.IMShow(src)

	dw := gocv.NewWindow("dst")
	defer dw.Close()

	trackBar1 := dw.CreateTrackbar("H_min",179)
	trackBar1.SetMin(50)
	//trackBar1.SetPos(80)

	trackBar2:= dw.CreateTrackbar("H_max",179)
	trackBar2.SetMin(80)

	for {
		hMin := trackBar1.GetPos()
		hMax := trackBar2.GetPos()
		gocv.InRangeWithScalar(srcHsv,gocv.Scalar{Val1: float64(hMin),Val2: 150,Val3: 0},gocv.Scalar{Val1: float64(hMax),Val2: 255,Val3: 255},&dst)
		dw.IMShow(dst)
		if dw.WaitKey(3) == 27 {
			break
		}
	}

	//gocv.WaitKey(0)
}
