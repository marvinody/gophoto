package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"bitbucket.org/marvinody/gophoto/crystallize"
	"bitbucket.org/marvinody/gophoto/edge"
	"bitbucket.org/marvinody/gophoto/glasstile"
	"bitbucket.org/marvinody/gophoto/maxrgb"
	"bitbucket.org/marvinody/gophoto/pixelize"
	"bitbucket.org/marvinody/gophoto/predator"
)

type Filter interface {
	Apply(image.Image) image.Image
}

func main() {

	effects := make([]Filter, 0, 10)

	if len(os.Args) > 1 {
		filterString := os.Args[1]
		filterStrings := strings.Split(filterString, ";")
		fmt.Fprintf(os.Stderr, "%v\n", filterStrings)
		for _, s := range filterStrings {
			s2 := strings.Split(s, ",")
			s = s2[0]
			s2 = s2[1:]
			fmt.Fprintf(os.Stderr, "%s\n", s)
			switch s {
			case "maxrgb":
				effect := maxrgb.Effect{}
				effects = append(effects, effect)
			case "edge":
				effect := edge.Effect{}
				effects = append(effects, effect)
			case "pixelize":
				effect := pixelize.Effect{}
				effects = append(effects, effect)
			case "predator":
				effect := predator.Effect{}
				effects = append(effects, effect)
			case "crystallize":
				effect := crystallize.Effect{20}
				effects = append(effects, effect)
			case "glasstile":
				effect := glasstile.Effect{20, 20}
				effects = append(effects, effect)
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	img, _, err := image.Decode(reader)
	if err != nil {
		file, _ := os.Open("shion.png")
		img, _, _ = image.Decode(file)
	}
	out := applyEffectChain(img,
		effects...,
	)
	png.Encode(os.Stdout, out)
}

func applyEffectChain(img image.Image, filters ...Filter) image.Image {
	temp := img
	for _, effect := range filters {
		temp = effect.Apply(temp)
	}
	return temp
}
