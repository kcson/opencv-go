package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	src := gocv.IMRead("candies.png", gocv.IMReadColor)
	if src.Empty() {
		fmt.Println("Img file read fail !!")
		os.Exit(2)
	}

	fmt.Println("src.shape : ", src.Size())
	fmt.Println("src.channel : ", src.Channels())
	fmt.Println("src.type : ", src.Type())

	srcHsv := gocv.NewMat()
	gocv.CvtColor(src, &srcHsv, gocv.ColorBGRToYCrCb)

	planes := gocv.Split(srcHsv)
	gocv.NewWindow("src").IMShow(src)
	gocv.NewWindow("planes[0]").IMShow(planes[0])
	gocv.NewWindow("planes[1]").IMShow(planes[1])
	gocv.NewWindow("planes[2]").IMShow(planes[2])

	gocv.WaitKey(0)

}
