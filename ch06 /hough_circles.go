package main

import (
    "fmt"
    "gocv.io/x/gocv"
    "image"
    "image/color"
    "os"
)

func main() {
    //src := gocv.IMRead("imgs/dial.jpg", gocv.IMReadColor)
    src := gocv.IMRead("imgs/coins1.jpg", gocv.IMReadColor)
    if src.Empty() {
        fmt.Println("image read fail!!")
        os.Exit(-1)
    }
    defer src.Close()

    gray := gocv.NewMat()
    defer gray.Close()
    gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)
    gocv.GaussianBlur(gray, &gray, image.Pt(0, 0), 1.0, 0, gocv.BorderDefault)

    srcWin := gocv.NewWindow("src")
    defer srcWin.Close()
    srcWin.IMShow(src)

    dstWin := gocv.NewWindow("dst")
    defer dstWin.Close()

    minRTrackBar := dstWin.CreateTrackbar("minRadius", 100)
    maxRTrackBar := dstWin.CreateTrackbar("maxRadius", 150)
    thresholdTrackBar := dstWin.CreateTrackbar("threshold", 100)
    minRTrackBar.SetPos(10)
    maxRTrackBar.SetPos(80)
    thresholdTrackBar.SetPos(40)

    circles := gocv.NewMat()
    defer circles.Close()
    dst := gocv.NewMat()
    defer dst.Close()

    for {
        rMin := minRTrackBar.GetPos()
        rMax := maxRTrackBar.GetPos()
        th := thresholdTrackBar.GetPos()

        gocv.HoughCirclesWithParams(gray, &circles, gocv.HoughGradient, 1, 50, 120, float64(th), rMin, rMax)
        src.CopyTo(&dst)
        if !circles.Empty() {
            for i := 0; i < circles.Cols(); i++ {
                v := circles.GetVecfAt(0, i)
                //fmt.Println(int(v[0]), int(v[1]))
                gocv.Circle(&dst, image.Pt(int(v[0]), int(v[1])), int(v[2]), color.RGBA{R: 255, G: 0, B: 0}, 2)
            }
        }
        dstWin.IMShow(dst)
        if gocv.WaitKey(3) == 27 {
            break
        }
    }
}
