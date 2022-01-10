package main

import (
    "fmt"
    "gocv.io/x/gocv"
    "os"
    "sync"
)

func main() {
    //coreNum := runtime.NumCPU()
    src := gocv.IMRead("lenna.bmp", gocv.IMReadGrayScale)
    if src.Empty() {
        fmt.Println("Image read fail !!")
        os.Exit(2)
    }
    defer src.Close()
    dst := gocv.NewMat()
    defer dst.Close()

    //var alpha uint16 = 2
    var alpha float32 = 1

    //src.ConvertTo(&dst, gocv.MatTypeCV64F)
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()
        src.ConvertToWithParams(&dst, gocv.MatTypeCV8U, 1+alpha, -128*alpha)
    }()
    go func() {
        defer wg.Done()
        src.ConvertToWithParams(&dst, gocv.MatTypeCV8U, 1+alpha, -128*alpha)
    }()
    wg.Wait()

    //f64bytes, _ := dst.DataPtrFloat64()
    //f64bytes, _ := dst.DataPtrUint16()
/*
    var wg sync.WaitGroup
    dataNum := len(f64bytes) / coreNum
    remainDataNum := len(f64bytes) % coreNum
    for i := 0; i < coreNum; i++ {
        wg.Add(1)
        startIndex := dataNum * i
        endIndex := dataNum * (i + 1)
        if i == coreNum-1 {
            endIndex += remainDataNum
        }
        go func(sIndex, eIndex int) {
            defer wg.Done()
            for j := sIndex; j < eIndex; j++ {
                f64bytes[j] = f64bytes[j]*(1+alpha) - 128*alpha
            }
        }(startIndex, endIndex)
    }
    wg.Wait()
/*
 */
    //for i, v := range f64bytes {
    //    f64bytes[i] = v*(1+alpha) - 128*alpha
    //}
    //floats.AddScaled(f64bytes, alpha, f64bytes)
    //floats.AddConst(-128*alpha, floats.AddScaledTo(f64bytes, make([]float64, len(f64bytes)), 1+alpha, f64bytes))

    //scaleBytes := make([]float64, len(f64bytes))
    //floats.AddConst(-128*alpha, scaleBytes)
    //floats.AddScaledTo(f64bytes, scaleBytes, 1+alpha, f64bytes)

    //if err != nil {
    //	fmt.Println(err.Error())
    //}
    //dense := mat.NewDense(src.Rows(), src.Cols(), f64bytes)
    ////dense.Add(dense, dense)
    //maxDense := mat.Max(dense)
    //fmt.Println("maxDense : ", maxDense)

    //gocv.Normalize(src, &dst, 0, 255, gocv.NormMinMax)

    //dst.ConvertTo(&dst, gocv.MatTypeCV8U)
    gocv.NewWindow("src").IMShow(src)
    gocv.NewWindow("dst").IMShow(dst)
    gocv.WaitKey(0)
}
