package main

import (
	"fmt"
	"gocv.io/x/gocv"
)

func main() {
	cap1, err := gocv.VideoCaptureFile("video1.mp4")
	if err != nil {
		fmt.Println("video capture fail !!")
		panic(err)
	}
	defer cap1.Close()
	cap2, err := gocv.VideoCaptureFile("video2.mp4")
	if err != nil {
		fmt.Println("video capture fail !!")
		panic(err)
	}
	defer cap2.Close()

	frameCnt1 := int(cap1.Get(gocv.VideoCaptureFrameCount))
	frameCnt2 := int(cap2.Get(gocv.VideoCaptureFrameCount))
	fps1 := cap1.Get(gocv.VideoCaptureFPS)
	fps2 := cap2.Get(gocv.VideoCaptureFPS)
	effectFrames := int(fps1 * 2)

	fmt.Println("frameCnt1 : ", frameCnt1)
	fmt.Println("frameCnt2 : ", frameCnt2)
	fmt.Println("fps1 : ", fps1)
	fmt.Println("fps2 : ", fps2)

	delay := int(1000 / fps1)

	w := cap1.Get(gocv.VideoCaptureFrameWidth)
	h := cap1.Get(gocv.VideoCaptureFrameHeight)
	out, _ := gocv.VideoWriterFile("output.avi", "DIVX", fps1, int(w), int(h), true)
	defer out.Close()

	window := gocv.NewWindow("output")
	defer window.Close()

	img1 := gocv.NewMat()
	defer img1.Close()

	img2 := gocv.NewMat()
	defer img2.Close()

	for i := 0; i < (frameCnt1 - effectFrames); i++ {
		if !cap1.Read(&img1) {
			break
		}
		out.Write(img1)

		window.IMShow(img1)
		window.WaitKey(delay)
	}
	fmt.Println(img1.Size())
	fmt.Println(img1.Channels())

	for i := 0; i < effectFrames; i++ {
		if !cap1.Read(&img1) || !cap2.Read(&img2) {
			break
		}

		dx := int(w) / effectFrames * i
		if dx == 0 {
			continue
		}
		temp1 := img1.ColRange(dx, int(w))
		temp2 := img2.ColRange(0, dx)
		gocv.Hconcat(temp2, temp1, &temp2)

		out.Write(temp2)

		window.IMShow(temp2)
		window.WaitKey(delay)
	}

	for i := effectFrames; i < frameCnt2; i++ {
		if !cap2.Read(&img2) {
			break
		}
		out.Write(img2)

		window.IMShow(img2)
		window.WaitKey(delay)
	}
}
