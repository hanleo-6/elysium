package render

import (
	"github.com/hanleo-6/elysium/internal/physics"
	"github.com/hanleo-6/elysium/internal/sim"
)

func (c *Canvas) DrawECSlug(slug *sim.ElysiaChlorotica) {
	color := SlugColor(slug.Fullness)
	// draw rough body shape
	for i, point := range slug.Skelly.Points {
		c.DrawCircleFill(int(point.Pos.X), int(point.Pos.Y), slug.BodyProfile[i], color)
	}
	// draw leaf wings
	leftWingPoints := make([]physics.Vec2, 0, 4)
	rightWingPoints := make([]physics.Vec2, 0, 4)

	maxWingSpan := 4.0

	for i := 1; i <= len(slug.BodyProfile)-3; i++ {
		p := slug.Skelly.Points[i].Pos
		pPrev := slug.Skelly.Points[i-1].Pos

		dir := p.Sub(pPrev)
		dist := dir.Length()
		if dist > 0 {
			dir = dir.Scale(1.0 / dist)
			norm := physics.Vec2{X: -dir.Y, Y: dir.X}

			offset := float64(slug.BodyProfile[i]) + (slug.WingsOpenFactor() * maxWingSpan)
			leftP := p.Add(norm.Scale(offset))
			rightP := p.Sub(norm.Scale(offset))

			leftWingPoints = append(leftWingPoints, leftP)
			rightWingPoints = append(rightWingPoints, rightP)
		}
	}
	for i := 1; i < len(leftWingPoints); i++ {
		c.DrawLine(int(leftWingPoints[i].X), int(leftWingPoints[i].Y),
		           int(leftWingPoints[i-1].X), int(leftWingPoints[i-1].Y), color)
	}

	for i := 1; i < len(rightWingPoints); i++ {
		c.DrawLine(int(rightWingPoints[i].X), int(rightWingPoints[i].Y),
		           int(rightWingPoints[i-1].X), int(rightWingPoints[i-1].Y), color)
	}

	// draw rhinopores on head
	head := slug.Head().Pos
	neck := slug.Skelly.Points[1].Pos

	dir := head.Sub(neck)
	dist := dir.Length()
	if dist > 0 {
		dir = dir.Scale(1.0 / dist)
		norm := physics.Vec2{X: -dir.Y, Y: dir.X}

		// move to constants
		rhinoporeLen := 3.0
		rhinoporeSpread := 2.0

		leftTip := head.Add(dir.Scale(rhinoporeLen)).Add(norm.Scale(rhinoporeSpread))
		rightTip := head.Add(dir.Scale(rhinoporeLen)).Sub(norm.Scale(rhinoporeSpread))

		c.DrawLine(int(head.X), int(head.Y), int(leftTip.X), int(leftTip.Y), color)
		c.DrawLine(int(head.X), int(head.Y), int(rightTip.X), int(rightTip.Y), color)
	}
}

