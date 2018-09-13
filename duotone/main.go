package duotone

import (
	"image"
	"image/color"
	"math"
)

type Effect struct {
	lowColor, highColor color.Color
}

func (e Effect) SetLowColor(clr color.Color) Effect {
	e.lowColor = clr
	return e
}
func (e Effect) SetHighColor(clr color.Color) Effect {
	e.highColor = clr
	return e
}

const gradientSteps int = 256

func lerp(a, b int, t float64) float64 {
	return float64(a)*t + float64(b)*(1-t)
}

func (e Effect) createGradientArray() []color.Color {
	arr := make([]color.Color, gradientSteps)
	for idx := range arr {
		l := float64(idx) / float64(gradientSteps-1)
		lC, rC := e.lowColor, e.highColor
		lR, lG, lB, _ := lC.RGBA()
		rR, rG, rB, _ := rC.RGBA()
		arr[idx] = color.RGBA{
			uint8(lerp(int(lR>>8), int(rR>>8), l)),
			uint8(lerp(int(lG>>8), int(rG>>8), l)),
			uint8(lerp(int(lB>>8), int(rB>>8), l)),
			0,
		}
	}
	return arr
}

func max(args ...float64) float64 {
	m := -math.MaxFloat64
	for _, el := range args {
		m = math.Max(0, el)
	}
	return m
}
func min(args ...float64) float64 {
	m := math.MaxFloat64
	for _, el := range args {
		m = math.Min(0, el)
	}
	return m
}

type HSL struct {
	H, S, L float64
}

func RGB2HSL(clr color.RGBA) HSL {
	r, g, b := float64(clr.R)/255, float64(clr.G)/255, float64(clr.B)/255
	max, min := max(r, g, b), min(r, g, b)
	h, s, l := 0.0, 0.0, (max+min)/2
	if max == min {
		h, s = 0, 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}
		rTern := 0.0
		if g < b {
			rTern = 6
		}
		switch max {
		case r:
			h = (g-b)/d + rTern
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h = h / 6
	}
	return HSL{h, s, l}
}

func (e Effect) Apply(img image.Image) image.Image {

	size := img.Bounds().Size()
	out := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	gradient := e.createGradientArray()
	for y := 0; y < size.Y; y += 1 {
		for x := 0; x < size.X; x += 1 {
			clr := img.At(x, y)
			hsl := RGB2HSL(color.RGBAModel.Convert(clr).(color.RGBA))
			idx := int(hsl.L * float64(gradientSteps))
			gradClr := gradient[idx]
			_, _, _, origA := clr.RGBA()
			gradR, gradG, gradB, _ := gradClr.RGBA()
			newClr := color.RGBA{
				uint8(gradR >> 8),
				uint8(gradG >> 8),
				uint8(gradB >> 8),
				uint8(origA >> 8),
			}

			out.Set(x, y, newClr)

		}

	}
	return out
}
