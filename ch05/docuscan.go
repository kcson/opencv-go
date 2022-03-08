package main

import (
    "fmt"
    "gocv.io/x/gocv"
    "image"
    "image/color"
    "os"
)

func drawROI(img gocv.Mat, corners [][]float32) gocv.Mat {
    cpy := gocv.NewMat()
    defer cpy.Close()
    img.CopyTo(&cpy)

    c1 := color.RGBA{B: 192, G: 192, R: 255}
    c2 := color.RGBA{B: 128, G: 128, R: 255}

    for _, pt := range corners {
        gocv.CircleWithParams(&cpy, image.Pt(int(pt[0]), int(pt[1])), 25, c1, -1, gocv.LineAA, 0)
    }
    gocv.Line(&cpy, image.Pt(int(corners[0][0]), int(corners[0][1])), image.Pt(int(corners[1][0]), int(corners[1][1])), c2, 2)
    gocv.Line(&cpy, image.Pt(int(corners[1][0]), int(corners[1][1])), image.Pt(int(corners[2][0]), int(corners[2][1])), c2, 2)
    gocv.Line(&cpy, image.Pt(int(corners[2][0]), int(corners[2][1])), image.Pt(int(corners[3][0]), int(corners[3][1])), c2, 2)
    gocv.Line(&cpy, image.Pt(int(corners[3][0]), int(corners[3][1])), image.Pt(int(corners[0][0]), int(corners[0][1])), c2, 2)

    disp := gocv.NewMat()
    gocv.AddWeighted(img, 0.3, cpy, 0.7, 0, &disp)

    return disp
}

func main() {
    src := gocv.IMRead("imgs/scanned.jpg", gocv.IMReadColor)
    if src.Empty() {
        fmt.Println("image read fail !!")
        os.Exit(1)
    }
    defer src.Close()

    h, w := src.Rows(), src.Cols()
    //dw := 500
    //dh := math.Round(float64(dw) * 297 / 210)

    srcQuad := [][]float32{{30, 30}, {30, float32(h) - 30}, {float32(w) - 30, float32(h) - 30}, {float32(w) - 30, 30}}
    //dstQuad := [][]float32{{0, 0}, {0, float32(dh) - 1}, {float32(dw) - 1, float32(dh) - 1}, {float32(dw) - 1, 0}}
    //dragSrc := []bool{false, false, false, false}

    //srcW := gocv.NewWindow("src")
    //defer srcW.Close()
    //srcW.IMShow(src)

    disp := drawROI(src, srcQuad)
    defer disp.Close()
    dispW := gocv.NewWindow("disp")
    defer dispW.Close()
    dispW.IMShow(disp)
}
