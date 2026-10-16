package sim

import (
	"math"
	"math/rand/v2"
)

type Saltmarsh struct {
	LightOriginX int
	LightOriginY int
	InitTheta    float64
	SpreadTheta  float64

	centreTheta  float64
	baseSpread   float64
	phases       [3]float64
	drift        float64
}

const pull = 0.15
const noise = 0.8

func NewSaltmarsh() *Saltmarsh {
	s := &Saltmarsh{centreTheta: 47.5, baseSpread: 24.0}
	for i := range s.phases {
		s.phases[i] = rand.Float64() * 2 * math.Pi
	}
	s.MoveLightband(0, 0)
	return s
}

func (s *Saltmarsh) MoveLightband(t, dt float64) {
	sway := 3.6 * math.Sin(0.31 * t + s.phases[0]) +
			1.5 * math.Sin(0.83 * t + s.phases[1]) +
			0.4 * math.Sin(2.17 * t + s.phases[2])

	// Pull back the light toward centre
	// 'Ornstein Uhlenbeck'
	// 'Euler-Maruyama' method
	// fancier the name, the better
	s.drift += -pull * s.drift * dt + noise * math.Sqrt(dt) * rand.NormFloat64()

	// Modulate width by a factor to spread and shrink it
	s.SpreadTheta = s.baseSpread * (1 + 0.25 * math.Sin(0.48 * t + s.phases[1]))

	// Sway centre
	centre := s.centreTheta + sway + s.drift
	// Derive InitTheta from centre
	s.InitTheta = centre - s.SpreadTheta / 2
}

