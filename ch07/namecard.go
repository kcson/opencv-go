package main

import "C"
import (
    "gocv.io/x/gocv"
    "image"
    "image/color"
    "math"
    "os"
    "sort"
)

func main() {
    //fileName := "imgs/namecard1.jpg"
    fileName := "imgs/vehicle.jpeg"
    if len(os.Args) > 1 {
        fileName = os.Args[1]
    }
    src := gocv.IMRead(fileName, gocv.IMReadColor)
    if src.Empty() {
        println("image read fail!!")
        os.Exit(1)
    }
    defer src.Close()

    //dw, dh := 800, 400
    //srcQuad := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {0, 0}, {0, 0}, {0, 0}})
    //dstQuad := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {0, dh}, {dw, dh}, {dw, 0}})

    gray := gocv.NewMat()
    defer gray.Close()
    gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)

    srcBin := gocv.NewMat()
    defer srcBin.Close()
    gocv.GaussianBlur(gray, &srcBin, image.Point{X: 0, Y: 0}, 2, 0, gocv.BorderConstant)
    gocv.AdaptiveThreshold(srcBin, &srcBin, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinaryInv, 19, 4)
    //gocv.Threshold(srcBin, &srcBin, 0, 255, gocv.ThresholdBinaryInv|gocv.ThresholdOtsu)
    //gocv.Threshold(gray, &srcBin, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)
    srcBinWin := gocv.NewWindow("gray")
    defer srcBinWin.Close()
    srcBinWin.IMShow(srcBin)

    //contours := gocv.FindContours(srcBin, gocv.RetrievalExternal, gocv.ChainApproxNone)
    contours := gocv.FindContours(srcBin, gocv.RetrievalList, gocv.ChainApproxNone)
    cpy := gocv.NewMat()
    defer cpy.Close()
    src.CopyTo(&cpy)
    for i := 0; i < contours.Size(); i++ {
        pts := contours.At(i)
        area := gocv.ContourArea(pts)
        if area < 500 {
            continue
        }
        rect := gocv.BoundingRect(pts)
        cnt := getConnectedComponents(srcBin, rect)
        println(cnt)
        if cnt < 7 || cnt > 20 {
            continue
        }
        gocv.Rectangle(&cpy, rect, color.RGBA{R: 255}, 2)
        /*approx := gocv.ApproxPolyDP(pts, gocv.ArcLength(pts, true)*0.02, true)
          if approx.Size() != 4 || !IsContourConvex(approx) {
          	continue
          }
          tempQuad := reorderPts(approx)
        */

        //srcQuad = tempQuad
        //gocv.Polylines(&cpy, gocv.NewPointsVectorFromPoints([][]image.Point{approx.ToPoints()}), true, color.RGBA{R: 255}, 3)*/
    }
    //pers := gocv.GetPerspectiveTransform(srcQuad, dstQuad)
    //dst := gocv.NewMat()
    //defer dst.Close()
    //gocv.WarpPerspective(src, &dst, pers, image.Pt(dw, dh))
    //
    //dstGray := gocv.NewMat()
    //defer dstGray.Close()
    //gocv.CvtColor(dst, &dstGray, gocv.ColorBGRToGray)
    //
    //dstWin := gocv.NewWindow("dst")
    //defer dstWin.Close()
    //
    //dstWin.IMShow(dstGray)
    srcWin := gocv.NewWindow("src")
    defer srcWin.Close()
    srcWin.IMShow(cpy)

    gocv.WaitKey(0)
}

func reorderPts(pts gocv.PointVector) gocv.PointVector {
    points := pts.ToPoints()
    sort.Slice(points, func(i, j int) bool {
        return points[i].X < points[j].X
    })
    if points[0].Y > points[1].Y {
        points[0], points[1] = points[1], points[0]
    }
    if points[2].Y < points[3].Y {
        points[2], points[3] = points[3], points[2]
    }
    return gocv.NewPointVectorFromPoints(points)
}

func IsContourConvex(curve gocv.PointVector) bool {
    if curve.Size() != 4 {
        return false
    }
    p1, p2, p3, p4 := curve.At(0), curve.At(1), curve.At(2), curve.At(3)
    if insideInTri(p1, p2, p3, p4) {
        return false
    }
    if insideInTri(p2, p3, p4, p1) {
        return false
    }
    if insideInTri(p3, p4, p1, p2) {
        return false
    }
    if insideInTri(p4, p1, p2, p3) {
        return false
    }

    return true
}

func insideInTri(p1, p2, p3, p image.Point) bool {
    a1 := gocv.ContourArea(gocv.NewPointVectorFromPoints([]image.Point{p1, p2, p}))
    a2 := gocv.ContourArea(gocv.NewPointVectorFromPoints([]image.Point{p2, p3, p}))
    a3 := gocv.ContourArea(gocv.NewPointVectorFromPoints([]image.Point{p3, p1, p}))
    a := gocv.ContourArea(gocv.NewPointVectorFromPoints([]image.Point{p1, p2, p3}))
    return math.Abs((a1+a2+a3)-a) < 0.1
}

func getConnectedComponents(mat gocv.Mat, pts image.Rectangle) int {
    mat = mat.Region(pts)
    defer mat.Close()
    gocv.MorphologyEx(mat, &mat, gocv.MorphOpen, gocv.NewMat())
    contours := gocv.FindContours(mat, gocv.RetrievalList, gocv.ChainApproxNone)
    if contours.Size() == 0 {
        return 0
    }
    componentCount := 0
    gocv.CvtColor(mat, &mat, gocv.ColorGrayToBGR)
    for i := 1; i < contours.Size(); i++ {
        pts := contours.At(i)
        area := gocv.ContourArea(pts)
        if area < 1000 {
            continue
        }
        if !guessVehicleNo(pts) {
           continue
        }
        componentCount++
        gocv.Rectangle(&mat, gocv.BoundingRect(pts), color.RGBA{R: 255}, 1)
    }
    if componentCount == 0 {
        return 0
    }
    matWin := gocv.NewWindow("mat")
    defer matWin.Close()
    matWin.IMShow(mat)
    gocv.WaitKey(0)

    return componentCount
}

func guessVehicleNo(pts gocv.PointVector) bool {
    rect := gocv.BoundingRect(pts)
    ratio := float32(rect.Dy()) / float32(rect.Dx())
    if ratio < 1.2 || ratio > 2.8 {
        return false
    }
    return true
}
