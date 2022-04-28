package main

import "C"
import (
    "github.com/otiai10/gosseract/v2"
    "gocv.io/x/gocv"
    "image"
    "image/color"
    "math"
    "os"
    "sort"
    "unicode"
    "unicode/utf8"
)

func main() {
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

    gray := gocv.NewMat()
    defer gray.Close()
    gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)

    srcBin := gocv.NewMat()
    defer srcBin.Close()
    gocv.GaussianBlur(gray, &srcBin, image.Point{X: 0, Y: 0}, 2, 0, gocv.BorderDefault)
    gocv.AdaptiveThreshold(srcBin, &srcBin, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinaryInv, 19, 4)
    srcBinWin := gocv.NewWindow("gray")
    defer srcBinWin.Close()
    srcBinWin.IMShow(srcBin)

    contours := gocv.FindContours(srcBin, gocv.RetrievalList, gocv.ChainApproxNone)
    cpy := gocv.NewMat()
    defer cpy.Close()
    src.CopyTo(&cpy)
    var minRect image.Rectangle
    for i := 0; i < contours.Size(); i++ {
        pts := contours.At(i)
        area := gocv.ContourArea(pts)
        if area < 100*100 {
            continue
        }
        rect := gocv.BoundingRect(pts)
        cnt := getConnectedComponents(srcBin, rect)
        println(cnt)
        if cnt < 4 || cnt > 20 {
            continue
        }
        if minRect.Empty() {
            minRect = rect
            continue
        }
        if (rect.Dx() * rect.Dy()) < (minRect.Dx() * minRect.Dy()) {
            minRect = rect
        }
    }
    //minRect.Min = image.Pt(minRect.Min.X+20, minRect.Min.Y)
    //minRect.Max = image.Pt(minRect.Max.X-20, minRect.Max.Y)
    gocv.Rectangle(&cpy, minRect, color.RGBA{R: 255}, 4)
    //dst := srcBin.Region(minRect)
    dst := cpy.Region(minRect)
    defer dst.Close()
    //dst = transformImage(dst)
    vehicleNo := getOcrText(dst)
    println(vehicleNo)

    srcWin := gocv.NewWindow("src")
    defer srcWin.Close()
    srcWin.IMShow(cpy)

    dstWin := gocv.NewWindow("dst")
    defer dstWin.Close()
    if !dst.Empty() {
        dstWin.IMShow(dst)
    }

    gocv.WaitKey(0)
}

func getOcrText(mat gocv.Mat) string {
    temp := gocv.NewMat()
    defer temp.Close()
    mat.CopyTo(&temp)
    gocv.CvtColor(temp, &temp, gocv.ColorBGRToGray)
    gocv.Threshold(temp, &temp, 0, 255, gocv.ThresholdOtsu|gocv.ThresholdBinaryInv)
    //gocv.MorphologyEx(mat, &mat, gocv.MorphDilate, gocv.NewMat())
    //gocv.MorphologyEx(mat, &mat, gocv.MorphDilate, gocv.NewMat())
    //gocv.MorphologyEx(mat, &mat, gocv.MorphDilate, gocv.NewMat())
    //gocv.MorphologyEx(mat, &mat, gocv.MorphDilate, gocv.NewMat())
    //gocv.MorphologyEx(mat, &mat, gocv.MorphDilate, gocv.NewMat())
    contours := gocv.FindContours(temp, gocv.RetrievalTree, gocv.ChainApproxNone)
    gocv.GaussianBlur(temp, &temp, image.Pt(0, 0), 1, 4, gocv.BorderDefault)
    gocv.Threshold(temp, &temp, 0, 255, gocv.ThresholdOtsu|gocv.ThresholdBinaryInv)
    result := ""
    for i := 0; i < contours.Size(); i++ {
        pts := contours.At(i)
        area := gocv.ContourArea(pts)
        if area < 50*50 {
            continue
        }
        rect := gocv.BoundingRect(pts)
        if !guessVehicleNo(pts) {
            continue
        }
        gocv.Rectangle(&mat, rect, color.RGBA{R: 255}, 4)
        result = getTextFromMat(temp.Region(rect))
        println(result)
    }
    //tempWin := gocv.NewWindow("temp")
    //defer tempWin.Close()
    //tempWin.IMShow(temp)

    result = getTextFromMat(temp)
    return result
}

func getTextFromMat(mat gocv.Mat) string {
    client := gosseract.NewClient()
    defer client.Close()
    text := ""
    gocv.IMWrite("vehicleNo.jpg", mat)
    //err := client.SetImageFromBytes(mat.ToBytes())
    err := client.SetImage("./vehicleNo.jpg")
    if err != nil {
        println(err.Error())
        return ""
    }
    err = client.SetPageSegMode(gosseract.PSM_SINGLE_WORD)
    if err != nil {
        println(err.Error())
        return ""
    }
    err = client.SetLanguage("kor")
    if err != nil {
        println(err.Error())
        return ""
    }
    text, err = client.Text()
    if err != nil {
        println(err.Error())
        return ""
    }
    if len(text) == 0 {
        return ""
    }
    println(text)
    result := ""
    b := []byte(text)
    for i := 0; i < len(b); {
        r, size := utf8.DecodeRune(b[i:])
        if (unicode.Is(unicode.Hangul, r) && r >= '가' && r <= '힣') || unicode.IsNumber(r) {
            result = result + string(r)
        }
        i += size
    }

    return result
}

func transformImage(dst gocv.Mat) gocv.Mat {
    //gocv.CopyMakeBorder()
    gocv.CvtColor(dst, &dst, gocv.ColorBGRToGray)
    gocv.Threshold(dst, &dst, 0, 255, gocv.ThresholdOtsu|gocv.ThresholdBinary)
    dw, dh := dst.Cols(), dst.Cols()/3
    srcQuad := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {0, 0}, {0, 0}, {0, 0}})
    dstQuad := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {0, dh}, {dw, dh}, {dw, 0}})
    contours := gocv.FindContours(dst, gocv.RetrievalList, gocv.ChainApproxNone)
    maxApprox := contours.At(0)
    for i := 0; i < contours.Size(); i++ {
        pts := contours.At(i)
        approx := gocv.ApproxPolyDP(pts, gocv.ArcLength(pts, true)*0.02, true)
        if approx.Size() != 4 || !IsContourConvex(approx) {
            continue
        }
        area := gocv.ContourArea(approx)
        if area < 300 {
            continue
        }
        if area > gocv.ContourArea(maxApprox) {
            maxApprox = approx
        }
    }
    srcQuad = reorderPts(maxApprox)
    pers := gocv.GetPerspectiveTransform(srcQuad, dstQuad)
    gocv.WarpPerspective(dst, &dst, pers, image.Pt(dw, dh))

    return dst
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
    //gocv.MorphologyEx(mat, &mat, gocv.MorphOpen, gocv.NewMat())
    contours := gocv.FindContours(mat, gocv.RetrievalList, gocv.ChainApproxNone)
    if contours.Size() == 0 {
        return 0
    }
    componentCount := 0
    gocv.CvtColor(mat, &mat, gocv.ColorGrayToBGR)
    for i := 1; i < contours.Size(); i++ {
        pts := contours.At(i)
        area := gocv.ContourArea(pts)
        if area < 200 {
            continue
        }
        if !guessVehicleNo(pts) {
            continue
        }
        similarCount := similarContour(contours, pts)
        if similarCount < 3 {
            continue
        }
        componentCount++
        gocv.Rectangle(&mat, gocv.BoundingRect(pts), color.RGBA{R: 255}, 1)
    }
    if componentCount == 0 {
        return 0
    }
    //matWin := gocv.NewWindow("mat")
    //defer matWin.Close()
    //matWin.IMShow(mat)
    //gocv.WaitKey(0)

    return componentCount
}

func similarContour(contours gocv.PointsVector, pts gocv.PointVector) int {
    similarCount := 0
    p1 := gocv.BoundingRect(pts)
    for i := 0; i < contours.Size(); i++ {
        p2 := gocv.BoundingRect(contours.At(i))
        areaDiff := math.Abs(float64(p1.Dx()*p1.Dy())-float64(p2.Dx()*p2.Dy())) / float64(p1.Dx()*p1.Dy())
        widthDiff := math.Abs(float64(p1.Dx())-float64(p2.Dx())) / float64(p1.Dx())
        heightDiff := math.Abs(float64(p1.Dy())-float64(p2.Dy())) / float64(p1.Dy())

        if areaDiff > 0.6 {
            continue
        }
        if widthDiff > 0.8 {
            continue
        }
        if heightDiff > 0.2 {
            continue
        }
        diagonalLength := math.Sqrt(float64(p1.Dx()*p1.Dx() + p1.Dy()*p1.Dy()))
        distance := math.Sqrt(float64((p1.Min.X-p2.Min.X)*(p1.Min.X-p2.Min.X) + (p1.Min.Y-p2.Min.Y)*(p1.Min.Y-p2.Min.Y)))
        if distance > diagonalLength*5 {
            continue
        }
        angleDiff := 0.
        dx := math.Abs(float64(p1.Min.X - p2.Min.X))
        dy := math.Abs(float64(p1.Min.Y - p2.Min.Y))
        if dx == 0 {
            angleDiff = 90
        } else {
            angleDiff = math.Atan(dy/dx) * (180 / math.Pi)
        }
        if angleDiff > 20.0 {
            continue
        }

        similarCount++
    }
    return similarCount
}

func guessVehicleNo(pts gocv.PointVector) bool {
    rect := gocv.BoundingRect(pts)
    ratio := float32(rect.Dy()) / float32(rect.Dx())
    if ratio < 1.0 || ratio > 4.0 {
        return false
    }
    return true
}
