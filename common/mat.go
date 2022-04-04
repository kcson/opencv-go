package common

import "gocv.io/x/gocv"

func NewMat(row, col int, data [][]float32, dataType gocv.MatType) *gocv.Mat {
    mat := gocv.Zeros(row, col, dataType)

    for i := 0; i < row; i++ {
        for j := 0; j < col; j++ {
            mat.SetFloatAt(i, j, data[i][j])
        }
    }

    return &mat
}
