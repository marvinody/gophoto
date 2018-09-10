package predator

import (
	"image"

	"bitbucket.org/marvinody/gophoto/edge"
	"bitbucket.org/marvinody/gophoto/maxrgb"
	"bitbucket.org/marvinody/gophoto/pixelize"
)

type Effect struct {
	BlurWidth, BlurHeight int
}

func (p Effect) Apply(img image.Image) image.Image {

	edgeEffect := edge.Effect{}
	maxrgbEffect := maxrgb.Effect{}
	pixelizeEffect := pixelize.Effect{}

	edgeTemp := edgeEffect.Apply(img)
	maxrgbTemp := maxrgbEffect.Apply(edgeTemp)
	pixelizeTemp := pixelizeEffect.Apply(maxrgbTemp)
	return pixelizeTemp

}
