package sim

type Saltmarsh struct {
	LightOriginX int
	LightOriginY int
}

func NewSaltmarsh() *Saltmarsh {
	return &Saltmarsh{
		LightOriginX: 0,
		LightOriginY: 0,
	}
}

func (s *Saltmarsh) UpdatePositions() {

}

// getGradient: returns a gradient that changes smoothly based on previous gradient
func (s *Saltmarsh) getGradient() {

}

