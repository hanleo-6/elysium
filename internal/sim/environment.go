package sim

import (
	"math"

	"github.com/hanleo-6/elysium/internal/physics"
)

// Environment encapsulates wakeGrid, saltmarsh,  TODO: algae, and slugs.
type Environment struct {
	grid *wakeGrid
	Saltmarsh *Saltmarsh
	Time float64
}

func NewEnvironment(dotWidth, dotHeight, cellSize int) *Environment {
	return &Environment{
		grid: newwakeGrid(dotWidth, dotHeight, cellSize),
		Saltmarsh: NewSaltmarsh(),
		Time: 0.0,
	}
}

// potential generates noise for a given point, to be used as a base for Curl noise
func (e *Environment) potential(x, y float64) float64 {
	// A_1 * sin(k_1 * x + w_1 * t) * cos(k_2 * y + w_2 * t) + A_2 * sin(k_3 * y - w_3 * t)
	// random constants
	return 15.0 * math.Sin(0.03*x + 0.4*e.Time) * math.Cos(0.04*y + 0.2*e.Time) + 1.0 * math.Sin(0.01*y - 0.3*e.Time)
}

// sampleCurl calculates Curl noise for a given point
//     potential on all points generates a given topological height map,
//     then curl noise is essentially vectors going around each height ring
func (e *Environment) sampleCurl(pos physics.Vec2) physics.Vec2{
	eps := 0.1
	a1 := e.potential(pos.X, pos.Y+eps)
	a2 := e.potential(pos.X, pos.Y-eps)
	dPsiDy := (a1 - a2) / (2 * eps)

	b1 := e.potential(pos.X+eps, pos.Y)
	b2 := e.potential(pos.X-eps, pos.Y)
	dPsiDx := (b1 - b2) / (2 * eps)

	return physics.Vec2{X: dPsiDy, Y: -dPsiDx}
}

func (e *Environment) flowAt(pos physics.Vec2) physics.Vec2 {
	// combine curl noise + slug wake for aquatic flow feeling
	curl := e.sampleCurl(pos)
	wake := e.grid.sampleWake(pos)
	return curl.Add(wake.Scale(wakeFactor))
}

// Tick updates all dynamic elements of Environment
// TODO: move algae to be encapsulated in Environment, move flow logic to saltmarsh, potentially move slugs to Env too
func (e *Environment) Tick(dt float64, algae []*Alga, slugs []*ElysiaChlorotica) {
	e.grid.tick(decayFactor)
	for _, alga := range algae {
		displacement := (e.flowAt(alga.Pos)).Scale(dt * displacementFactor)
		alga.Pos = alga.Pos.Add(displacement)
	}
	for _, slug := range slugs {
		pos := slug.Head().Pos
		vel := pos.Sub(slug.Head().PrevPos)
		e.grid.injectForce(pos, vel)
	}
	e.Saltmarsh.MoveLightband()
	e.Time += dt
 }

