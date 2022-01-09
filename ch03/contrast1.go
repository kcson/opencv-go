package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("lenna.bmp", gocv.IMReadGrayScale)
	if src.Empty() {
		fmt.Println("Image read fail !!")
		os.Exit(2)
	}
	defer src.Close()
	dst := gocv.NewMat()
	defer dst.Close()

	//var alpha uint16 = 2
	var alpha float64 = 2

	src.ConvertTo(&dst, gocv.MatTypeCV64F)
	f64bytes, _ := dst.DataPtrFloat64()
	//f64bytes, _ := dst.DataPtrUint16()
	for i, v := range f64bytes {
		f64bytes[i] = v*(1+alpha) - 128*alpha
	}
	//floats.AddScaled(f64bytes, alpha, f64bytes)
	//floats.AddConst(-128*alpha, floats.AddScaledTo(f64bytes, make([]float64, len(f64bytes)), 1+alpha, f64bytes))

	//scaleBytes := make([]float64, len(f64bytes))
	//floats.AddConst(-128*alpha, scaleBytes)
	//floats.AddScaledTo(f64bytes, scaleBytes, 1+alpha, f64bytes)

	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//dense := mat.NewDense(src.Rows(), src.Cols(), f64bytes)
	////dense.Add(dense, dense)
	//maxDense := mat.Max(dense)
	//fmt.Println("maxDense : ", maxDense)

	//gocv.Normalize(src, &dst, 0, 255, gocv.NormMinMax)

	dst.ConvertTo(&dst, gocv.MatTypeCV8U)
	gocv.NewWindow("src").IMShow(src)
	gocv.NewWindow("dst").IMShow(dst)
	gocv.WaitKey(0)
}
