package sim

import (
	"math"

	"github.com/hanleo-6/elysium/internal/physics"
)

type Stage int
const (
	StageJuvenile Stage = iota
	StageAdult
)

type ElysiaChlorotica struct {
	Skelly          *physics.Skeleton
	Fullness        float64
	BodyProfile     []int
	Stage           Stage
}

const jointNumEC int = 7
const nomRange float64 = 1.0

func NewECSlug(initialX, initialY int, segmentLength float64) *ElysiaChlorotica {
	return &ElysiaChlorotica{
		Skelly: physics.NewSpine(physics.Vec2{X: float64(initialX), Y: float64(initialY)}, jointNumEC, segmentLength),
		Fullness: 0.0,
		BodyProfile: []int{
			2,
			2,
			3,
			2,
			1,
			1,
			0,
		},
	}
}

func (s *ElysiaChlorotica) Head() *physics.Point {
	return s.Skelly.Points[0]
}

func (s *ElysiaChlorotica) Tick(algae []*Alga, dt float64) {
	// if s.Fullness >= 1.0 {
	// 	// seek light zone + loaf
	// } else {
	// 	s.SeekFood(algae, dt)
	// }
	s.SeekFood(algae, dt)
}

func (s *ElysiaChlorotica) WingsOpenFactor() float64 {
	if s.Fullness < 1 {
		return 0.0
	}
	return s.Fullness
}

// TODO: consider concurrent seeking
func (s *ElysiaChlorotica) SeekFood(algae []*Alga, dt float64) {
	target, dist := s.acquireTarget(algae)
	if dist < nomRange {
		target.Eat()
		s.Fullness += target.NomFactor()
	}

	accel := s.Skelly.Seek(target.Pos, dt)
	s.Skelly.Step(dt, accel)
}

func (s *ElysiaChlorotica) acquireTarget(algae []*Alga) (*Alga, float64) {
	// naive check
	var closest *Alga
	minDist := 999999999.0
	for _, alga := range algae {
		if alga.Eaten {
			continue
		}
		offset := alga.Pos.Sub(s.Head().Pos)
		dist := math.Sqrt(offset.X*offset.X + offset.Y*offset.Y)
		if dist < minDist {
			closest = alga
			minDist = dist
		}
	}
	return closest, minDist
}

