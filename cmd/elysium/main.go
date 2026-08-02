package main

import(
	"fmt"
	"github.com/hanleo-6/elysium/internal/render"
)

func main() {
	c := render.New(25, 25)

	c.DrawCircle(c.DotWidth()/3, c.DotHeight()/2, 3.0, render.ColorAdult)
	c.DrawLine(0, 0, 5, 2, render.ColorDeepWater)

	fmt.Print(c.Render())
}

