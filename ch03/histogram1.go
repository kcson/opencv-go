package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
)

func main() {
	src := gocv.IMRead("lenna.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("Img read fail!!")
		os.Exit(2)
	}

	hist := gocv.NewMat()
	gocv.CalcHist([]gocv.Mat{src}, []int{0}, gocv.NewMat(), &hist, []int{256}, []float64{0, 256}, false)

	gocv.NewWindow("src").IMShow(src)
	gocv.WaitKey(1)
	gocv.NewWindow("hist").IMShow(drawHistogramV2(hist))
	gocv.WaitKey(0)
}

func drawHistogramV2(hist gocv.Mat) gocv.Mat {
	histImage := gocv.NewMatWithSizeFromScalar(gocv.Scalar{Val1: 255}, 200, 256, gocv.MatTypeCV8U)
	_, histMax, _, _ := gocv.MinMaxIdx(hist)
	for x := 0; x < 256; x++ {
		p1 := image.Point{X: x, Y: 200}
		p2 := image.Point{X: x, Y: int(200 - hist.GetFloatAt(x, 0)*200/histMax)}
		gocv.Line(&histImage, p1, p2, color.RGBA{R: 0, G: 0, B: 0}, 1)
	}

	return histImage
}

func drawHistogram(hist gocv.Mat) gocv.Mat {

	dHist := gocv.NewMat()
	defer dHist.Close()

	//set matrix size to be shown
	histW := 512
	histH := 512
	size := hist.Size()[0]
	binW := int(float64(histW) / float64(size))
	histImage := gocv.NewMatWithSize(512, 400, gocv.MatTypeCV8U)

	gocv.Normalize(hist, &dHist, 0, float64(histImage.Rows()), gocv.NormMinMax) // normalize to show easily

	for idx := 1; idx < size; idx++ {
		gocv.Line(
			&histImage,
			image.Point{X: binW * (idx - 1), Y: histH - int(dHist.GetFloatAt(idx-1, 0))},
			image.Point{X: binW * (idx), Y: histH - int(dHist.GetFloatAt(idx, 0))},
			color.RGBA{R: 255, G: 255, B: 255},
			2)
	}

	return histImage
}
