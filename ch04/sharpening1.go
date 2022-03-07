package main

import (
    "fmt"
    "gocv.io/x/gocv"
    "os"
)

func main() {
    src := gocv.IMRead("imgs/rose.bmp",gocv.IMReadGrayScale)
    if src.Empty() {
        fmt.Println("image read fail!!")
        os.Exit(1)
    }
    defer src.Close()

    blur := gocv.NewMat()
    defer blur.Close()

    srcWindow := gocv.NewWindow("src")
    defer srcWindow.Close()

    srcWindow.IMShow(src)

    gocv.WaitKey(0 )
}
