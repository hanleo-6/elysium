package physics

import (
	"math"
)

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vec2) Scale(k float64) Vec2 {
	return Vec2{
		X: v.X * k,
		Y: v.Y * k,
	}
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.X * v.X + v.Y * v.Y)
}

type Point struct {
	Pos     Vec2
	PrevPos Vec2
}

type Constraint struct {
	P1, P2     *Point
	RestLength float64
}

