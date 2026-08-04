package sim

import (
	"math/rand"

	"github.com/hanleo-6/elysium/internal/physics"
)

type Alga struct {
	Pos physics.Vec2
	Radius float64 // doubles as 'nom factor', i.e. how much a slug grows by eating it
	Eaten bool
}

func NewAlga(x, y int) *Alga {
	return &Alga{
		Pos: physics.Vec2{X: float64(x), Y: float64(y)},
		Radius: 1.0 + rand.Float64(),
		Eaten: false,
	}
}

func (a *Alga) Eat() {
	a.Eaten = true
}

func (a *Alga) NomFactor() float64 {
	return a.Radius * 0.1
}

