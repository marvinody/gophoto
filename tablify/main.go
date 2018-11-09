package tablify

import (
	"image"
	"image/color"
	"math"
)

type Effect struct {
	BG  color.Color
	Tol float64
}

func (e Effect) Apply(img image.Image) image.Image {

	bnds := img.Bounds()
	oldWidth, oldHeight := bnds.Dx(), bnds.Dy()
	newWidth, newHeight := oldWidth*2-1, oldHeight*2-1
	out := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	rects := Tablify(img, e.Tol)
	for y := 0; y < newHeight; y += 1 {
		for x := 0; x < newWidth; x += 1 {
			// I don't like this but set bg on all pixels on new first
			out.Set(x, y, e.BG)
		}
	}

	for _, line := range rects {
		for _, rect := range line {
			for x := rect.x; x < rect.x+rect.width; x += 1 {
				for y := rect.y; y < rect.y+rect.height; y += 1 {
					out.Set(2*x, 2*y, rect.color)
					canExpandX := x+1 < rect.x+rect.width
					canExpandY := y+1 < rect.y+rect.height
					if canExpandX {
						out.Set(2*x+1, 2*y, rect.color)
					}
					if canExpandY {
						out.Set(2*x, 2*y+1, rect.color)
					}
					if canExpandX && canExpandY {
						out.Set(2*x+1, 2*y+1, rect.color)
					}

				}
			}
		}
	}
	return out
}

func colorDist(A, B color.Color) float64 {
	aR, aG, aB, aA := A.RGBA()
	bR, bG, bB, bA := B.RGBA()
	return math.Sqrt(
		math.Pow(float64(aR>>8)-float64(bR>>8), 2) +
			math.Pow(float64(aG>>8)-float64(bG>>8), 2) +
			math.Pow(float64(aB>>8)-float64(bB>>8), 2) +
			math.Pow(float64(aA>>8)-float64(bA>>8), 2),
	)
}

type Rect struct {
	x, y          int
	width, height int
	color         color.Color
}

func boolify(img image.Image) [][]bool {
	min, max := img.Bounds().Min, img.Bounds().Max
	bools := make([][]bool, max.Y-min.Y)
	for idx := range bools {
		bools[idx] = make([]bool, max.X-min.X)
	}
	return bools
}

func Tablify(img image.Image, tol float64) [][]Rect {
	min, max := img.Bounds().Min, img.Bounds().Max
	rects := make([][]Rect, 0, max.Y-min.Y)
	bools := boolify(img)
	for y := min.Y; y < max.Y; y += 1 {
		rowRect := make([]Rect, 0)
		for x := min.X; x < max.X; x += 1 {
			if bools[y][x] {
				//if we've already 'scanned' it, just keep moving
				continue
			}
			rect := rectify(img, bools, x, y, tol)

			rowRect = append(rowRect, rect)
		}
		rects = append(rects, rowRect)

	}
	return rects
}

func rectify(img image.Image, bools [][]bool, x, y int, tolerance float64) Rect {
	r := Rect{x, y, 1, 1, img.At(x, y)}
	_, max := img.Bounds().Min, img.Bounds().Max
	bools[y][x] = true
	/* we can expand IF,
	within Bounds
	pixel not used
	pixel same color
	*/
	canExpand := func(plusX, plusY int) bool {
		newX, newY := x+(r.width-1)+plusX,
			y+(r.height-1)+plusY
		return newX < max.X && newY < max.Y && // within bounds
			!bools[newY][newX] && // not used
			colorDist(r.color, img.At(newX, newY)) <= tolerance // "same" color
	}
	for canExpand(1, 0) {
		// just set the flag now
		bools[y+r.height-1][x+r.width] = true
		r.width += 1 // if we can expand right, just do it
	}

	for canExpand(0, 1) {
		rowGood := true
		for i := 0; i < r.width; i += 1 {
			rowGood = rowGood && canExpand(-i, 1)
		}
		if rowGood {
			for i := 0; i < r.width; i++ {
				bools[y+r.height][x+i] = true
			}
			r.height += 1
		} else {
			break
		}
	}
	return r
}
