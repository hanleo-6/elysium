package render

import (
	"math"

	"github.com/hanleo-6/elysium/internal/sim"
)

func (c *Canvas) DrawLightband(s *sim.Saltmarsh) {
	initialRad := s.InitTheta * math.Pi / 180.0
	spreadRad := s.SpreadTheta * math.Pi / 180.0
	c.DrawCone(s.LightOriginX, s.LightOriginY, initialRad, spreadRad, ColorLightBand)
}

