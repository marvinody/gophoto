package pixelize

import (
	"image"
	"image/color"
	"math"
)

type Effect struct {
	BlurWidth, BlurHeight int
}

const defaultBlurWidth int = 10
const defaultBlurHeight int = 10

func blockIndex(pos, size int) int {
	if pos < 0 {
		return (pos+1)/size - 1
	} else {
		return pos / size
	}
}

func meanRectangle(img image.Image, rect image.Rectangle) color.Color {
	clrTotal := make([]float64, 4)
	rectX, rectY := rect.Bounds().Min.X, rect.Bounds().Min.Y
	for y := rectY; y < rectY+rect.Dy(); y += 1 {
		for x := rectX; x < rectX+rect.Dx(); x += 1 {
			clr := img.At(x, y)
			r, g, b, a := clr.RGBA()
			clrTotal[0] += float64(r & 0xff)
			clrTotal[1] += float64(g & 0xff)
			clrTotal[2] += float64(b & 0xff)
			clrTotal[3] += float64(a & 0xff)
		}
	}
	pixels := rect.Dx() * rect.Dy()
	clrTotal[0] = math.Round(clrTotal[0] / float64(pixels))
	clrTotal[1] = math.Round(clrTotal[1] / float64(pixels))
	clrTotal[2] = math.Round(clrTotal[2] / float64(pixels))
	clrTotal[3] = math.Round(clrTotal[3] / float64(pixels))
	return color.RGBA{
		uint8(clrTotal[0]),
		uint8(clrTotal[1]),
		uint8(clrTotal[2]),
		uint8(clrTotal[3]),
	}
}
func setRectangle(out *image.RGBA, rect image.Rectangle, clr color.Color) {
	rectX, rectY := rect.Bounds().Min.X, rect.Bounds().Min.Y
	for y := rectY; y < rectY+rect.Dy(); y += 1 {
		for x := rectX; x < rectX+rect.Dx(); x += 1 {
			out.Set(x, y, clr)
		}
	}
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
func (p Effect) Apply(img image.Image) image.Image {

	out := image.NewRGBA(img.Bounds())

	sizeX, sizeY := p.BlurWidth, p.BlurHeight
	if sizeX <= 0 {
		sizeX = defaultBlurWidth
	}
	if sizeY <= 0 {
		sizeY = defaultBlurHeight
	}
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	for y := 0; y < height; y += sizeY {
		for x := 0; x < width; x += sizeX {
			rectEndX, rectEndY := clamp(x+sizeX, 0, width), clamp(y+sizeY, 0, height)
			rect := image.Rect(x, y, rectEndX, rectEndY)
			clr := meanRectangle(img, rect)
			setRectangle(out, rect, clr)
		}

	}
	return out
}
