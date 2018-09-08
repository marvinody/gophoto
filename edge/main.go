package edge

import (
	"image"
	"image/color"
	"math"
)

type Effect struct {
	Amt float64
}

func clamp(t, min, max int) int {
	if t < min {
		return min
	}
	if t > max {
		return max
	}
	return t
}

var vKernel [9]int64 = [9]int64{
	-1, 0, +1,
	-2, 0, +2,
	-1, 0, +1,
}
var hKernel [9]int64 = [9]int64{
	-1, -2, -1,
	+0, +0, +0,
	+1, +2, +1,
}

func applyKernel3x3(pixels []color.RGBA, kernel []int64) []int64 {
	temp := make([]int64, 3)
	// R
	for i := 0; i < len(pixels); i += 1 {
		temp[0] += int64(pixels[i].R) * kernel[i]
		temp[1] += int64(pixels[i].G) * kernel[i]
		temp[2] += int64(pixels[i].B) * kernel[i]
	}

	return temp
}

func sqrt(iterations int, n float64) float64 {
	guess := 6.0 * math.Pow(10, n)
	for i := 0; i < iterations; i += 1 {
		guess = (guess + n/guess) / 2
	}
	return guess
}

// length of pr... should be imageWidth + 2
func edgeSobel(row int, pr, cr, nr []color.RGBA, amount float64) []color.Color {
	imgWidth := len(pr) - 2
	out := make([]color.Color, imgWidth)
	for i := 1; i < len(pr)-1; i += 1 {
		pxlBuffer := make([]color.RGBA, 0, 9)
		pxlBuffer = append(pxlBuffer, pr[i-1:i+2]...)
		pxlBuffer = append(pxlBuffer, cr[i-1:i+2]...)
		pxlBuffer = append(pxlBuffer, nr[i-1:i+2]...)
		vGrad := applyKernel3x3(pxlBuffer, vKernel[:])
		hGrad := applyKernel3x3(pxlBuffer, hKernel[:])
		vMagSqR, hMagSqR := math.Pow(float64(vGrad[0]), 2)*amount, math.Pow(float64(hGrad[0]), 2)*amount
		vMagSqG, hMagSqG := math.Pow(float64(vGrad[1]), 2)*amount, math.Pow(float64(hGrad[1]), 2)*amount
		vMagSqB, hMagSqB := math.Pow(float64(vGrad[2]), 2)*amount, math.Pow(float64(hGrad[2]), 2)*amount
		// so I lost about an hour on this shit
		// because golang think it's super clever to do a modulus on uint8
		// something like z := 256, x = uint8(z) will give you 0
		// good luck thinking it's the cast that does this shit
		// literally the worst part of go so far
		// I spent idk how many times going over the algo
		// debugging line by line, seeing which part messed up
		// and no matter how many times I went through it
		// it always got messed up
		// I had everything more condensed then so it made me think
		// the bug had to be in my portion of the code (because I'm not amazing)
		// but nope. piece of shit idea someone had and put it in
		r := clamp(int(math.Sqrt(vMagSqR+hMagSqR)), 0, 255)
		g := clamp(int(math.Sqrt(vMagSqG+hMagSqG)), 0, 255)
		b := clamp(int(math.Sqrt(vMagSqB+hMagSqB)), 0, 255)

		a := cr[i].A
		out[i-1] = color.RGBA{uint8(r), uint8(g), uint8(b), a}
	}
	return out
}

func createBufferedRGBARowWithExtra(img image.Image, row int) []color.RGBA {
	width := img.Bounds().Size().X
	colors := make([]color.RGBA, width+2)
	colors[0] = color.RGBAModel.Convert(img.At(0, row)).(color.RGBA)
	for x := 0; x < width; x += 1 {
		colors[x+1] = color.RGBAModel.Convert(img.At(x, row)).(color.RGBA)
	}
	colors[width+1] = color.RGBAModel.Convert(img.At(width-1, row)).(color.RGBA)
	return colors
}

func (e Effect) Apply(img image.Image) image.Image {
	if e.Amt <= 0 {
		e.Amt = 1
	}
	out := image.NewRGBA(img.Bounds())
	var pr, cr, nr []color.RGBA
	size := img.Bounds().Max
	pr = createBufferedRGBARowWithExtra(img, 0)
	cr = pr
	nr = createBufferedRGBARowWithExtra(img, 1)
	for y := 0; y < size.Y; y += 1 {
		pixelRow := edgeSobel(y, pr, cr, nr, e.Amt)
		for x, clr := range pixelRow {
			out.Set(x, y, clr)
		}
		pr = cr
		cr = nr
		nr = createBufferedRGBARowWithExtra(img, clamp(y+1, 0, size.Y))
	}
	return out
}
