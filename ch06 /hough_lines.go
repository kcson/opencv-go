package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/building.jpg", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("image read fail")
		os.Exit(-1)
	}
	defer src.Close()

	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	srcWin.IMShow(src)

	edges := gocv.NewMat()
	defer edges.Close()

	gocv.Canny(src, &edges, 50, 150)

	lines := gocv.NewMat()
	defer lines.Close()
	gocv.HoughLinesPWithParams(edges, &lines, 1, math.Pi/180, 160, 30, 5)
	gocv.CvtColor(edges, &edges, gocv.ColorGrayToBGR)

	//lines.ConvertTo(&lines, gocv.MatTypeCV32S)

	fmt.Println(lines.Rows())
	fmt.Println(lines.Cols())
	if !lines.Empty() {
		for i := 0; i < lines.Rows(); i++ {
			//fmt.Println(i)
			//fmt.Println(lines.GetIntAt(i, 0))
			//fmt.Println(lines.GetIntAt(i, 1))
			//fmt.Println(lines.GetIntAt(i, 2))
			//fmt.Println(lines.GetIntAt(i, 3))
			//fmt.Println(lines.GetVeciAt(i,0))
			v := lines.GetVeciAt(i,0)
			gocv.Line(
				&edges,
				image.Pt(int(v[0]), int(v[1])),
				image.Pt(int(v[2]), int(v[3])),
				color.RGBA{R: 255, G: 0, B: 0},
				2,
			)
		}
	}

	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()
	dstWin.IMShow(edges)

	gocv.WaitKey(0)
}
