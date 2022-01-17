package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	cap1, err := gocv.VideoCaptureFile("imgs/woman.mp4")
	if err != nil {
		fmt.Println("video open fail !!")
		os.Exit(1)
	}
	defer cap1.Close()
	cap2, err := gocv.VideoCaptureFile("imgs/raining.mp4")
	if err != nil {
		fmt.Println("video open fail !!")
		os.Exit(1)
	}
	defer cap2.Close()

	//frameCnt1 := cap1.Get(gocv.VideoCaptureFrameCount)
	//frameCnt2 := cap2.Get(gocv.VideoCaptureFrameCount)

	fps := cap1.Get(gocv.VideoCaptureFPS)
	delay := int(1000 / fps)

	w := cap1.Get(gocv.VideoCaptureFrameWidth)
	h := cap1.Get(gocv.VideoCaptureFrameHeight)
	out, _ := gocv.VideoWriterFile("output.avi", "DIVX", fps, int(w), int(h), true)
	defer out.Close()

	doComposit := false

	frame1 := gocv.NewMat()
	defer frame1.Close()
	frame2 := gocv.NewMat()
	defer frame2.Close()
	hsv := gocv.NewMat()
	defer hsv.Close()
	mask := gocv.NewMat()
	defer mask.Close()

	frame1Window := gocv.NewWindow("frame1")
	defer frame1Window.Close()
	for {
		if !cap1.Read(&frame1) {
			break
		}
		if doComposit {
			if !cap2.Read(&frame2) {
				break
			}
			gocv.CvtColor(frame1, &hsv, gocv.ColorBGRToHSV)
			gocv.InRangeWithScalar(hsv, gocv.Scalar{Val1: 50, Val2: 150, Val3: 0}, gocv.Scalar{Val1: 70, Val2: 255, Val3: 255}, &mask)
			frame2.CopyToWithMask(&frame1, mask)
		}
		frame1Window.IMShow(frame1)
		out.Write(frame1)
		key := gocv.WaitKey(delay)
		if key == 27 {
			break
		} else if key == 32 {
			doComposit = !doComposit
		}
	}
}
