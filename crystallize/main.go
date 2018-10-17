package crystallize

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"

	"bitbucket.org/marvinody/kdtree"
)

type voronoi struct {
	idx   int
	x, y  int
	color color.Color
}

type Effect struct {
	Cells int
}

// from https://bitbucket.org/marvinody/crystallize/src/master/
func (e Effect) Apply(img image.Image) image.Image {
	cells := e.Cells
	if cells < 1 {
		cells = 1
	}
	size := img.Bounds().Size()
	maxGridSpacing := int((size.X + size.Y) / 2 / cells)
	minGridSpacing := int(maxGridSpacing / 2)
	if maxGridSpacing < 2 {
		maxGridSpacing = 2
	}
	if minGridSpacing < 1 {
		minGridSpacing = 1
	}
	gridRange := maxGridSpacing - minGridSpacing
	split := (maxGridSpacing - minGridSpacing) / 2
	pick := func(x, y int) (int, int) {
		xr, yr := rand.Intn(gridRange), rand.Intn(gridRange)
		if xr > split {
			xr += gridRange
		}
		if yr > split {
			yr += gridRange
		}
		return x + xr, y + yr
	}
	tree := &kdtree.KDTree{}
	indToVoronoi := make(map[int]voronoi)
	// so iterate over and create the "regions" we want
	fmt.Fprintf(os.Stderr, "Starting Voronoi Partitioning\n")
	for y := 0; y < size.Y; y += maxGridSpacing {
		for x := 0; x < size.X; x += maxGridSpacing {
			// voronoi pixel for this region
			vx, vy := pick(x, y)
			if vx >= size.X {
				vx = size.X - 1
			}
			if vy >= size.Y {
				vy = size.Y - 1
			}
			idx, err := tree.Insert([]float64{float64(vx), float64(vy)})
			if err != nil {
				panic(err)
			}
			indToVoronoi[idx] = voronoi{idx: idx, x: vx, y: vy, color: img.At(vx, vy)}
		}
	}
	tree.Balance()

	fmt.Fprintf(os.Stderr, "Done, Starting Pixel Partitioning\n")
	crystal := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	pixels := size.X * size.Y
	pxlCnt := 0
	reminders := 0
	sizeToRemind := 25 // 1/sizeToRemind fractional updates
	for y := 0; y < size.Y; y += 1 {
		for x := 0; x < size.X; x += 1 {
			pxlCnt += 1
			if pxlCnt > (pixels / sizeToRemind * reminders) {
				fmt.Fprintf(os.Stderr, "%3.f%% done\n", float32(reminders)/float32(sizeToRemind)*100)
				reminders += 1
			}
			tuple := []float64{float64(x), float64(y)}
			nearest := tree.FindNearestNeighborInd(tuple)
			v := indToVoronoi[nearest]
			crystal.Set(x, y, v.color)
		}
	}

	return crystal

}
