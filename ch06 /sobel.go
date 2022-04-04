package main

import (
    "fmt"
    "gocv.io/x/gocv"
    "opencv-go/common"
    "os"
)

func main() {
    src := gocv.IMRead("imgs/lenna.bmp", gocv.IMReadGrayScale)
    if src.Empty() {
        fmt.Println("image read fail!!")
        os.Exit(-1)
    }
    defer src.Close()

    srcWin := gocv.NewWindow("src")
    defer srcWin.Close()
    srcWin.IMShow(src)

    kernel := common.NewMat(
        3,
        3,
        [][]float32{
            {-1, 0, 1},
            {-2, 0, 2},
            {-1, 0, 1},
        }, gocv.MatTypeCV32F)
    defer kernel.Close()

    dx := gocv.NewMat()
    defer dx.Close()
    dy := gocv.NewMat()
    defer dy.Close()

    //gocv.Filter2D(src, &dx, -1, *kernel, image.Pt(-1, -1), 128, gocv.BorderDefault)

    //sobel 함수
    gocv.Sobel(src, &dx, -1, 1, 0, 3, 1, 128, gocv.BorderDefault)
    gocv.Sobel(src, &dy, -1, 0, 1, 3, 1, 128, gocv.BorderDefault)

    dxWin := gocv.NewWindow("dx")
    defer dxWin.Close()
    dxWin.IMShow(dx)

    dyWin := gocv.NewWindow("dy")
    defer dyWin.Close()
    dyWin.IMShow(dy)

    gocv.WaitKey(0)
}
