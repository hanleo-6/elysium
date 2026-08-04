package render

import (
	"fmt"
	"math/rand"

	"github.com/lucasb-eyer/go-colorful"
)
// Just imagine wherever it says color it is spelled colour..

func mustHex(hex string) colorful.Color {
	c, err := colorful.Hex(hex)
	if err != nil {
		panic(fmt.Sprintf("render/palette: invalid hex colour %q: %v", hex, err))
	}
	return c
}

// Exported Palette
var (
	// Water column
	ColorDeepWater    = mustHex("#513b56")
	ColorShallowWater = mustHex("#525174")
	ColorLightBand    = mustHex("#ACABC4")

	// Elysia chlorotica
	ColorJuvenile = mustHex("#8C4A32")
	ColorAdult    = mustHex("#2E8B57")

	// Algae
	ColorAlgaeLight = mustHex("#AACC00")
	ColorAlgaeDark  = mustHex("#557A00")

	// Portal
	ColorPortal = mustHex("#aa3e98")
)

func SlugColor(fullness float64) colorful.Color {
	return ColorJuvenile.BlendHcl(ColorAdult, fullness).Clamped()
}

func AlgaColor() colorful.Color {
	return ColorAlgaeLight.BlendHcl(ColorAlgaeDark, rand.Float64()).Clamped()
}

func WaterColor(depth float64) colorful.Color {
	return ColorShallowWater.BlendLab(ColorDeepWater, depth).Clamped()
}

