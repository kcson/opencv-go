package common

import (
	"gocv.io/x/gocv"
)

func NewMat(row, col int, data [][]float32, dataType gocv.MatType) *gocv.Mat {
	mat := gocv.Zeros(row, col, dataType)

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			mat.SetFloatAt(i, j, data[i][j])
		}
	}

	return &mat
}

func Clip(mat gocv.Mat, min, max float32) gocv.Mat {
	for i := 0; i < mat.Rows(); i++ {
		for j := 0; j < mat.Cols(); j++ {
			if mat.GetFloatAt(i, j) < min {
				mat.SetFloatAt(i, j, min)
			}
			if mat.GetFloatAt(i, j) > max {
				mat.SetFloatAt(i, j, max)
			}
		}
	}
	return mat
}
