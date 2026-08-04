package physics

type Skeleton struct {
	Points []*Point
	Constraints []*Constraint
	Damping float64
	MaxSpeed float64
	MaxForce float64
}

const numIterations int = 5

// NewSpine creates a simplified skeleton of a line of points
func NewSpine(root Vec2, count int, segmentLength float64) *Skeleton {
	p := make([]*Point, count)
	c := make([]*Constraint, count-1)
	for i := 0; i < count; i++ {
		pos := Vec2{X: root.X + segmentLength * float64(i), Y: root.Y}
		p[i] = &Point{Pos: pos, PrevPos: pos}
		if i > 0 {
			c[i-1] = &Constraint{P1: p[i-1], P2: p[i], RestLength: segmentLength}
		}
	}
	return &Skeleton{
		Points: p,
		Constraints: c,
		Damping: 0.98,
		MaxSpeed: 42.0,
		MaxForce: 380.0,
	} // TODO: put these into config
}

// Step calculates new position for each point in the skeleton
func (s *Skeleton) Step(dt float64, headAccel Vec2) {
	// update position of head with acceleration
	if len(s.Points) > 0 {
		head := s.Points[0]
		vel := (head.Pos.Sub(head.PrevPos)).Scale(s.Damping)
		head.PrevPos = head.Pos
		head.Pos = head.Pos.Add(vel).Add(headAccel.Scale(dt * dt))
	}
	// update rest of points
	for i := 1; i < len(s.Points); i++ {
		p := s.Points[i]
		vel := ((p.Pos).Sub(p.PrevPos)).Scale(s.Damping)
		p.PrevPos = p.Pos
		p.Pos = (p.Pos).Add(vel)
	}

	// contract and stretch segments by their constraints
	for range numIterations {
		for _, c := range s.Constraints {
			delta := c.P2.Pos.Sub(c.P1.Pos)
			dist := delta.Length()
			if dist == 0 {
				continue
			}
			diff := (dist - c.RestLength) / dist
			c.P1.Pos = c.P1.Pos.Add(delta.Scale(0.5 * diff))
			c.P2.Pos = c.P2.Pos.Sub(delta.Scale(0.5 * diff))
		}
	}
}

// Seek adds force to direct the head of the skeleton towards a target
// https://www.red3d.com/cwr/steer/gdc99/
func (s *Skeleton) Seek(target Vec2, dt float64) Vec2 {
	head := s.Points[0]
	curOffset := target.Sub(head.Pos)
	dist := curOffset.Length()
	if dist == 0 {
		return Vec2{}
	}
	seekVel := curOffset.Scale(s.MaxSpeed / dist)
	curVel := head.Pos.Sub(head.PrevPos).Scale(1.0 / dt)

	steerForce := seekVel.Sub(curVel)

	if steerForce.Length() > s.MaxForce {
		steerForce = steerForce.Scale(s.MaxForce / steerForce.Length())
	}
	return steerForce
}

