package glasstile

import (
	"image"
)

type Effect struct {
	TileWidth, TileHeight int
}

func (e Effect) Apply(img image.Image) image.Image {
	/*
	  Algo taken almost verbatim from https://github.com/OpenCL/GEGL-OpenCL/blob/d5e8cd3f96c050faa96f8cded31c919780ce6f0a/operations/common/tile-glass.c#L2
	*/
	tileWidth, tileHeight := e.TileWidth, e.TileHeight
	if tileWidth < 2 {
		tileWidth = 2
	}
	if tileHeight < 2 {
		tileHeight = 2
	}
	size := img.Bounds().Size()
	dstSize := size
	glassed := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	clamp := func(t, min, max int) int {
		if t < min {
			return min
		}
		if t > max {
			return max
		}
		return t
	}

	x1, y1 := 0, 0
	y2 := y1 + dstSize.Y

	xHalf, yHalf := tileWidth/2, tileHeight/2
	xPlus, yPlus := tileWidth%2, tileHeight%2

	dstXOffs := x1 % (tileWidth + xPlus)

	//srcX0 := x1 - dstXOffs
	xRightAbyss := 2 * ((x1 + dstSize.X) % tileWidth)
	if xRightAbyss > tileWidth-2 {
		xRightAbyss = tileWidth - 2
	}
	//srcRowWidth := dstXOffs + size.X + xRightAbyss
	yOffs := y1 % tileHeight
	yMiddle := y1 - yOffs
	if yOffs >= yHalf {
		yMiddle += tileHeight
		yOffs -= tileHeight
	}
	// loop through rows
	for row := y1; row < y2; row += 1 {
		yPixel2 := yMiddle + yOffs*2
		yPixel2 = clamp(yPixel2, 0, size.Y-1)
		yOffs += 1
		if yOffs == yHalf {
			yMiddle += tileHeight
			yOffs = -(yHalf + yPlus)
		}

		xOffs := x1 % tileWidth
		xMiddle := x1 - xOffs
		if xOffs >= xHalf {
			xMiddle += tileWidth
			xOffs -= tileWidth
		}
		for col := 0; col < size.X; col += 1 {
			xPixel1 := (xMiddle + xOffs - x1)
			xPixel2 := 0
			if xMiddle+xOffs*2+dstXOffs < size.X {
				xPixel2 = xMiddle + xOffs*2 - x1 + dstXOffs
			} else {
				xPixel2 = xMiddle + xOffs - x1 + dstXOffs

			}

			clr := img.At(xPixel2, yPixel2)
			glassed.Set(xPixel1, row, clr)

			xOffs += 1
			if xOffs == xHalf {
				xMiddle += tileWidth
				xOffs = -(xHalf + xPlus)
			}
		}

	}
	return glassed
}
