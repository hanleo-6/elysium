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

