package main

import "C"
import (
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
	"os"
	"regexp"
	"sort"
	"unicode"
	"unicode/utf8"
)

type Chain interface {
	Execute(param *ChainParam)
	SetNext(chain Chain) Chain
}
type ChainParam struct {
	Src       gocv.Mat
	Pre       gocv.Mat
	Dst       gocv.Mat
	VehicleNo string
}

type PreProcess struct {
	Next Chain
}

func (pre *PreProcess) Execute(param *ChainParam) {
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(param.Src, &gray, gocv.ColorBGRToGray)

	srcBin := gocv.NewMat()
	gocv.GaussianBlur(gray, &srcBin, image.Pt(0, 0), 2, 0, gocv.BorderDefault)
	gocv.AdaptiveThreshold(srcBin, &srcBin, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinaryInv, 19, 4)
	param.Pre = srcBin

	srcBinWin := gocv.NewWindow("PreProcess")
	defer srcBinWin.Close()
	srcBinWin.IMShow(srcBin)

	if pre.Next != nil {
		pre.Next.Execute(param)
	}
}
func (pre *PreProcess) SetNext(chain Chain) Chain {
	pre.Next = chain
	return chain
}

type CutVehicleRegion struct {
	Next Chain
}

func (cvr *CutVehicleRegion) Execute(param *ChainParam) {
	contours := gocv.FindContours(param.Pre, gocv.RetrievalList, gocv.ChainApproxNone)
	var minRect image.Rectangle
	for i := 0; i < contours.Size(); i++ {
		pts := contours.At(i)
		rect := gocv.BoundingRect(pts)
		ratio := float32(rect.Dy()) / float32(rect.Dx())
		if ratio > 0.5 {
			continue
		}
		area := gocv.ContourArea(pts)
		if area < 80*80 {
			continue
		}
		//gocv.Rectangle(&param.Src, rect, color.RGBA{R: 255}, 3)
		cnt := getConnectedComponents(param.Pre, rect)
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
	if minRect.Empty() {
		param.Dst = gocv.NewMat()
	} else {
		minRect.Min.X = minRect.Min.X + 20
		minRect.Max.X = minRect.Max.X - 20
		param.Dst = param.Src.Region(minRect)
		gocv.Rectangle(&param.Src, minRect, color.RGBA{R: 255}, 3)
	}

	dstWin := gocv.NewWindow("CutVehicleRegion")
	defer dstWin.Close()
	dstWin.IMShow(param.Dst)
	srcBinWin := gocv.NewWindow("param.Src")
	defer srcBinWin.Close()
	srcBinWin.IMShow(param.Src)

	if cvr.Next != nil {
		cvr.Next.Execute(param)
	}
}
func (cvr *CutVehicleRegion) SetNext(chain Chain) Chain {
	cvr.Next = chain
	return chain
}

type GetVehicleNo struct {
	Next Chain
}

func (gvn *GetVehicleNo) Execute(param *ChainParam) {
	dilateCount := 0 //팽창
	erodeCount := 2 //침식
	gray := gocv.NewMat()
	defer gray.Close()
	bin := gocv.NewMat()
	defer bin.Close()

	gocv.CvtColor(param.Dst, &gray, gocv.ColorBGRToGray)
	gocv.GaussianBlur(gray, &bin, image.Pt(0, 0), 1, 4, gocv.BorderDefault)
	thres := gocv.Threshold(bin, &bin, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)
	//침식
	for i := 0; i < erodeCount; i++ {
		gocv.MorphologyEx(bin, &bin, gocv.MorphErode, gocv.NewMat())
	}
	// 팽창
	for i := 0; i < dilateCount; i++ {
		gocv.MorphologyEx(bin, &bin, gocv.MorphDilate, gocv.NewMat())
	}

	gocv.CopyMakeBorder(bin, &bin, 10, 10, 10, 10, gocv.BorderConstant, color.RGBA{R: 0, G: 0, B: 0})
	println("threshold : ", thres)

	r4n, _ := regexp.Compile("[0-9]{4}")
	r, _ := regexp.Compile("[0-9]{2,3}[가-힣]{1}[0-9]{4}|[가-힣]{2}[0-9]{2}[가-힣]{1}[0-9]{4}")
	vehicleNo := getTextFromMat(bin)
	temp := r.FindString(vehicleNo)
	if temp != "" {
		vehicleNo = temp
	} else if len(vehicleNo) >= 4 {
		temp = vehicleNo[len(vehicleNo)-4:]
		temp = r4n.FindString(temp)
		if temp != "" {
			vehicleNo = temp
		}
	}
	if temp == "" {
		for i := 0; i < 10; i++ {
			thres = thres + 10
			println("threshold : ", thres)
			gocv.GaussianBlur(gray, &bin, image.Pt(0, 0), 1, 4, gocv.BorderDefault)
			gocv.Threshold(bin, &bin, thres, 255, gocv.ThresholdBinary)
			//침식
			for i := 0; i < erodeCount; i++ {
				gocv.MorphologyEx(bin, &bin, gocv.MorphErode, gocv.NewMat())
			}
			// 팽창
			for i := 0; i < dilateCount; i++ {
				gocv.MorphologyEx(bin, &bin, gocv.MorphDilate, gocv.NewMat())
			}
			gocv.CopyMakeBorder(bin, &bin, 10, 10, 10, 10, gocv.BorderConstant, color.RGBA{R: 0, G: 0, B: 0})
			vehicleNo = getTextFromMat(bin)
			temp = r.FindString(vehicleNo)
			if temp != "" {
				vehicleNo = temp
				break
			}
			if len(vehicleNo) < 4 {
				continue
			}
			temp = vehicleNo[len(vehicleNo)-4:]
			println(temp)
			temp = r4n.FindString(temp)
			if temp != "" {
				vehicleNo = temp
				break
			}
		}
	}

	param.VehicleNo = vehicleNo
	targetWin := gocv.NewWindow("GetVehicleNo")
	defer targetWin.Close()
	targetWin.IMShow(bin)

	if gvn.Next != nil {
		gvn.Next.Execute(param)
	}
}
func (gvn *GetVehicleNo) SetNext(chain Chain) Chain {
	gvn.Next = chain
	return chain
}

type GetVehicleNoFromChar struct {
	Next Chain
}

func (gvn *GetVehicleNoFromChar) Execute(param *ChainParam) {
	gray := gocv.NewMat()
	defer gray.Close()
	bin := gocv.NewMat()
	defer bin.Close()

	gocv.CvtColor(param.Dst, &gray, gocv.ColorBGRToGray)
	gocv.GaussianBlur(gray, &bin, image.Pt(0, 0), 1, 4, gocv.BorderDefault)
	gocv.Threshold(bin, &bin, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)
	gocv.MorphologyEx(bin, &bin, gocv.MorphErode, gocv.NewMat())
	gocv.MorphologyEx(bin, &bin, gocv.MorphErode, gocv.NewMat())

	targetWin := gocv.NewWindow("RegionChar")
	defer targetWin.Close()
	contours := gocv.FindContours(bin, gocv.RetrievalList, gocv.ChainApproxNone)
	for i := 0; i < contours.Size(); i++ {
		pts := contours.At(i)
		rect := gocv.BoundingRect(pts)
		//gocv.Rectangle(&param.Dst, rect, color.RGBA{R: 255}, 2)
		mat := bin.Region(rect)
		targetWin.IMShow(mat)
		gocv.WaitKey(0)
		getTextFromMat(mat)
	}

	targetWin.IMShow(param.Dst)

	if gvn.Next != nil {
		gvn.Next.Execute(param)
	}
}
func (gvn *GetVehicleNoFromChar) SetNext(chain Chain) Chain {
	gvn.Next = chain
	return chain
}

type ReleaseResource struct {
	Next Chain
}

func (rr *ReleaseResource) Execute(param *ChainParam) {
	if param.Src.Ptr() != nil {
		param.Src.Close()
	}
	if param.Pre.Ptr() != nil {
		param.Pre.Close()
	}
	if param.Dst.Ptr() != nil {
		param.Dst.Close()
	}
	if rr.Next != nil {
		rr.Next.Execute(param)
	}
	gocv.WaitKey(0)
}
func (rr *ReleaseResource) SetNext(chain Chain) Chain {
	rr.Next = chain
	return chain
}

func main() {
	fileName := "imgs/vehicle16.jpeg"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}
	src := gocv.IMRead(fileName, gocv.IMReadColor)
	if src.Empty() {
		println("image read fail!!")
		os.Exit(1)
	}
	param := new(ChainParam)
	param.Src = src

	pre := new(PreProcess)
	pre.SetNext(new(CutVehicleRegion)).SetNext(new(GetVehicleNo)).SetNext(new(ReleaseResource))
	pre.Execute(param)

	println("==================")
	println(param.VehicleNo)
	println("==================")
	//gocv.WaitKey(0)
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
	err = client.SetPageSegMode(gosseract.PSM_SINGLE_LINE)
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
