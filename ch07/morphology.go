package main

import (
    "gocv.io/x/gocv"
    "image"
    "os"
)

func main() {
    src := gocv.IMRead("imgs/circuit.bmp", gocv.IMReadGrayScale)
    if src.Empty() {
        println("Image load fail!!")
        os.Exit(1)
    }
    defer src.Close()

    se := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 3))
    defer se.Close()

    dst1 := gocv.NewMat()
    defer dst1.Close()
    dst2 := gocv.NewMat()
    defer dst2.Close()

    gocv.Erode(src, &dst1, se)
    gocv.Dilate(src, &dst2, gocv.NewMat())

    srcWin := gocv.NewWindow("src")
    defer srcWin.Close()
    dst1Win := gocv.NewWindow("dst1")
    defer dst1Win.Close()
    dst2Win := gocv.NewWindow("dst2")
    defer dst2Win.Close()


    srcWin.IMShow(src)
    dst1Win.IMShow(dst1)
    dst2Win.IMShow(dst2)

    gocv.WaitKey(0)
}
