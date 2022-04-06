package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

func main() {
	src := gocv.IMRead("imgs/coins1.jpg", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("image read fail!!")
		os.Exit(1)
	}
	defer src.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)
	gocv.GaussianBlur(gray, &gray, image.Pt(0, 0), 1, 0, gocv.BorderDefault)

	circles := gocv.NewMat()
	defer circles.Close()
	gocv.HoughCirclesWithParams(gray, &circles, gocv.HoughGradient, 1, 50, 150, 40, 20, 80)

	sumOfMoney := 0
	dst := gocv.NewMat()
	defer dst.Close()

	var mask gocv.Mat
	defer mask.Close()

	hsv := gocv.NewMat()
	defer hsv.Close()

	src.CopyTo(&dst)
	if !circles.Empty() {
		for i := 0; i < circles.Cols(); i++ {
			v := circles.GetVecfAt(0, i)
			cx, cy, radius := v[0], v[1], v[2]
			gocv.Circle(&dst, image.Pt(int(cx), int(cy)), int(radius), color.RGBA{R: 255}, 2)

			x1, y1, x2, y2 := int(cx-radius), int(cy-radius), int(cx+radius), int(cy+radius)
			crop := dst.ColRange(x1, x2)
			crop = crop.RowRange(y1, y2)
			ch, cw := crop.Rows(), crop.Cols()
			mask = gocv.Zeros(ch, cw, gocv.MatTypeCV8U)
			gocv.Circle(&mask, image.Pt(ch/2, cw/2), int(radius), color.RGBA{R: 255, G: 255, B: 255}, -1)

			gocv.CvtColor(crop, &hsv, gocv.ColorBGRToHSV)
			hue := gocv.Split(hsv)[0]
			hue.AddUChar(40)
			meanOfHue := hue.MeanWithMask(mask).Val1

			won := 100
			if meanOfHue < 90 {
				won = 10
			}
			sumOfMoney += won

			gocv.PutText(&crop, fmt.Sprintf("%d", won), image.Pt(20, 50), gocv.FontHersheySimplex, 0.75, color.RGBA{R: 255}, 2)
		}
	}
	gocv.PutText(&dst, fmt.Sprintf("%d", sumOfMoney), image.Pt(20, 50), gocv.FontHersheySimplex, 0.75, color.RGBA{R: 255}, 2)
	srcWin := gocv.NewWindow("src")
	defer srcWin.Close()
	dstWin := gocv.NewWindow("dst")
	defer dstWin.Close()

	srcWin.IMShow(src)
	dstWin.IMShow(dst)

	gocv.WaitKey(0)
}
