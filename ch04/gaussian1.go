package main

import (
    "fmt"
    "gocv.io/x/gocv"
    "image"
    "os"
)

func main() {
    src := gocv.IMRead("imgs/rose.bmp", gocv.IMReadGrayScale)
    if src.Empty() {
        fmt.Println("img read fail!!")
        os.Exit(1)
    }
    defer src.Close()
    dst := gocv.NewMat()
    defer dst.Close()

    gocv.GaussianBlur(src, &dst, image.Point{X: 0, Y: 0}, 2.0, 0, gocv.BorderConstant)

    srcWindow := gocv.NewWindow("src")
    defer srcWindow.Close()
    srcWindow.IMShow(src)

    dstWindow := gocv.NewWindow("dst")
    defer dstWindow.Close()
    dstWindow.MoveWindow(src.Cols(), 0)
    dstWindow.IMShow(dst)

    gocv.WaitKey(0)
}
