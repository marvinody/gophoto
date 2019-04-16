package greyscale

import (
	"image"
	"image/color"
	"math"
)

const (
	DefaultRCoeff = 0.2126
	DefaultGCoeff = 0.7152
	DefaultBCoeff = 0.0722
)

const (
	floatEqualityDelta = 0.00001
)

type Effect struct {
	RCoeff, GCoeff, BCoeff float64
}

func (e Effect) greyscale(clr color.Color) color.Color {
	r, g, b, a := clr.RGBA()
	gR := float64(r&0xff) * e.RCoeff
	gG := float64(g&0xff) * e.GCoeff
	gB := float64(b&0xff) * e.BCoeff
	greyTotal := gR + gG + gB
	grey := uint8(greyTotal)
	return color.RGBA{grey, grey, grey, uint8(a & 255)}
}

func (e Effect) Apply(img image.Image) image.Image {

	if floatEquals(e.RCoeff, 0) && floatEquals(e.GCoeff, 0) && floatEquals(e.BCoeff, 0) {
		e.RCoeff = DefaultRCoeff
		e.GCoeff = DefaultGCoeff
		e.BCoeff = DefaultBCoeff
	}

	out := image.NewRGBA(img.Bounds())
	size := img.Bounds().Max

	for y := 0; y < size.Y; y += 1 {
		for x := 0; x < size.X; x += 1 {
			out.Set(x, y, e.greyscale(img.At(x, y)))
		}
	}
	return out
}

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) <= floatEqualityDelta
}
