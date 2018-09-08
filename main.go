package main

import (
	"bufio"
	"image"
	"image/png"
	"os"

	"bitbucket.org/marvinody/gophoto/edge"
	"bitbucket.org/marvinody/gophoto/maxrgb"
	"bitbucket.org/marvinody/gophoto/pixelize"
)

type filter interface {
	Apply(image.Image) image.Image
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	img, _, err := image.Decode(reader)
	if err != nil {
		file, _ := os.Open("Valve_original.PNG")
		img, _, _ = image.Decode(file)
	}
	out := applyEffectChain(img,
		pixelize.Effect{10, 10},
		maxrgb.Effect{},
		edge.Effect{1},
	)
	png.Encode(os.Stdout, out)
}

func applyEffectChain(img image.Image, filters ...filter) image.Image {
	temp := img
	for _, effect := range filters {
		temp = effect.Apply(temp)
	}
	return temp
}
