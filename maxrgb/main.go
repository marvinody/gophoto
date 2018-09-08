package maxrgb

import (
	"image"
	"image/color"
)

type Effect struct {
	amt float64
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

func maxRGB(clr color.Color) color.Color {
	var max uint32 = 0
	channelFlags := 0
	var threshold uint32 = 10
	r, g, b, a := clr.RGBA()
	flag := func(c uint32, b uint) {
		c = c & 255
		if c+threshold >= max && c-threshold <= max {
			channelFlags |= 1 << b
		} else if c > max {
			max = c
			channelFlags = 1 << b
		}
	}
	check := func(b uint) uint8 {
		if channelFlags&(1<<b) > 0 {
			return uint8(max)
		}
		return 0
	}
	flag(r, 0)
	flag(g, 1)
	flag(b, 2)
	rR := check(0)
	rG := check(1)
	rB := check(2)
	return color.RGBA{rR, rG, rB, uint8(a & 255)}
}

func (e Effect) Apply(img image.Image) image.Image {

	out := image.NewRGBA(img.Bounds())
	size := img.Bounds().Max

	for y := 0; y < size.Y; y += 1 {
		for x := 0; x < size.X; x += 1 {
			out.Set(x, y, maxRGB(img.At(x, y)))
		}
	}
	return out
}
