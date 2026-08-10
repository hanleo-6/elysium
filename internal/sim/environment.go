package sim

import (
	"math"

	"github.com/hanleo-6/elysium/internal/physics"
)

// wakeGrid is used to track the 'wakes' of slugs, which push algae away from them
type wakeGrid struct {
	Cols int
	Rows int
	// CellWidth int
	// CellHeight int
	CellSize int
	Cells []physics.Vec2
}

const wakeFactor float64 = 2.0 // weigh wake more heavily than normal flow
const decayFactor float64 = 0.98
const displacementFactor float64 = 1.8 // make algae move around a little more

func newwakeGrid(dotWidth, dotHeight, cellSize int) *wakeGrid {
	cols := dotWidth / cellSize
	rows := dotHeight / cellSize
	return &wakeGrid{
		Cols: cols,
		Rows: rows,
		CellSize: cellSize,
		Cells: make([]physics.Vec2, rows*cols),
	}
}

func (g *wakeGrid) center() physics.Vec2 {
	return physics.Vec2{X: float64(g.Rows*g.CellSize) / 2.0, Y: float64(g.Cols*g.CellSize) / 2.0}
}

// injectForce adds the force from the wake of a seaslug to the wakeGrid
func (g *wakeGrid) injectForce(pos physics.Vec2, force physics.Vec2) {
	col := int(pos.X) / g.CellSize
	row := int(pos.Y) / g.CellSize
	if col >= 0 && col < g.Cols && row >= 0 && row < g.Rows {
		ix := row * g.Cols + col
		g.Cells[ix] = g.Cells[ix].Add(force)
	}
}

// sampleWake does bilinear interpolation between adjacent wakeGrid cells to calculate force on given point
// (reminder that a wakeGrid cell is generally bigger than a terminal dot)
func (g *wakeGrid) sampleWake(pos physics.Vec2) physics.Vec2 {
	gy := pos.Y / float64(g.CellSize)
	gx := pos.X / float64(g.CellSize)
	i := int(math.Floor(gy))
	j := int(math.Floor(gx))
	ty := gy - float64(i)
	tx := gx - float64(j)

	var vTotal physics.Vec2
	// check bounds
	if i == 0 || i+1 >= g.Rows || j == 0 || j+1 >= g.Cols {
		vTotal = g.center().Sub(pos)
	} else {
		vTop := g.Cells[i*g.Cols+(j+1)].Scale(1-tx).Add(g.Cells[(i+1)*g.Cols+(j+1)].Scale(tx))
		vBottom := g.Cells[i*g.Cols+j].Scale(1-tx).Add(g.Cells[(i+1)*g.Cols+j].Scale(tx))
		vTotal = vTop.Scale(ty).Add(vBottom.Scale(1-ty))
	}
	return vTotal
}

func (g *wakeGrid) tick(decay float64) {
	for row := 0; row < g.Rows; row++ {
		for col := 0; col < g.Cols; col++ {
			ix := row * g.Cols + col
			g.Cells[ix] = g.Cells[ix].Scale(decay)
		}
	}
}

// Environment encapsulates wakeGrid, TODO: algae, and slugs.
type Environment struct {
	Grid *wakeGrid
	Time float64
}

func NewEnvironment(dotWidth, dotHeight, cellSize int) *Environment {
	return &Environment{
		Grid: newwakeGrid(dotWidth, dotHeight, cellSize),
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
	wake := e.Grid.sampleWake(pos)
	return curl.Add(wake.Scale(wakeFactor))
}

// Tick updates all dynamic elements of Environment
// TODO: move algae to be encapsulated in Environment, move flow logic to saltmarsh, potentially move slugs to Env too
func (e *Environment) Tick(dt float64, algae []*Alga, slugs []*ElysiaChlorotica) {
	e.Grid.tick(decayFactor)
	for _, alga := range algae {
		displacement := (e.flowAt(alga.Pos)).Scale(dt * displacementFactor)
		alga.Pos = alga.Pos.Add(displacement)
	}
	for _, slug := range slugs {
		pos := slug.Head().Pos
		vel := pos.Sub(slug.Head().PrevPos)
		e.Grid.injectForce(pos, vel)
	}
	e.Time += dt
 }

