package sim

import (
	"math/rand/v2"
)

type Saltmarsh struct {
	LightOriginX int
	LightOriginY int
	InitTheta    float64
	SpreadTheta  float64
	expandVel    float64
	moveVel      float64
}

const expandScalar float64 = 1.5;
const moveScalar float64 = 0.1;

func getRandScalar(scalar float64) float64 {
	return rand.Float64() * scalar
}

func getRandomVel() float64 {
	// sign := rand.Float64 > 0.5 ? 1 : -1
	// return sign * rand.Float64()
	return 1.0
}

func NewSaltmarsh() *Saltmarsh {
	return &Saltmarsh{
		LightOriginX: 0,
		LightOriginY: 0,
		InitTheta:    45.0,
		SpreadTheta:  20.0,
		expandVel:    getRandomVel(),
		moveVel:      getRandomVel(),
	}
}

func (s *Saltmarsh) MoveLightband() {
	s.expandVel = getRandomVel() * getRandScalar(expandScalar)
	s.moveVel = getRandomVel() * getRandScalar(moveScalar)

	s.InitTheta += s.moveVel
	s.SpreadTheta += s.expandVel
}

